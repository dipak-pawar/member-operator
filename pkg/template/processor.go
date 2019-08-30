package template

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"context"
	templatev1 "github.com/openshift/api/template/v1"
	"github.com/openshift/library-go/pkg/template/generator"
	"github.com/openshift/library-go/pkg/template/templateprocessing"
	errs "github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Processor struct {
	cl           client.Client
	scheme       *runtime.Scheme
	codecFactory serializer.CodecFactory
}

func NewProcessor(cl client.Client, scheme *runtime.Scheme) *Processor {
	return &Processor{cl: cl, scheme: scheme, codecFactory: serializer.NewCodecFactory(scheme)}
}

func (p Processor) ProcessAndApply(tmplContent []byte, values map[string]string) error {
	objs, err := p.Process(tmplContent, values)
	if err != nil {
		return errs.Wrap(err, "failed while processing template")
	}

	return apply(p.cl, objs)
}

func apply(cl client.Client, objs []runtime.RawExtension) error {
	// sorting before applying to maintain correct order
	sort.Sort(ByKind(objs))

	for _, rawObj := range objs {
		obj := rawObj.Object
		if obj == nil {
			continue
		}
		gvk := obj.GetObjectKind().GroupVersionKind()
		if err := createOrUpdateObj(cl, obj); err != nil {
			return errs.Wrapf(err, "unable to create resource of kind: %s, version: %s", gvk.Kind, gvk.Version)
		}

		//if err := execute(cl, obj, http.MethodPost); err != nil {
		//	return errs.Wrapf(err, "unable to create resource of kind: %s, version: %s", gvk.Kind, gvk.Version)
		//}
	}
	return nil
}

func createOrUpdateObj(cl client.Client, obj runtime.Object) error {
	if err := cl.Create(context.TODO(), obj); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errs.Wrapf(err, "failed to create object %v", obj)
		}

		if err = cl.Update(context.TODO(), obj); err != nil {
			return errs.Wrapf(err, "failed to update object %v", obj)
		}
		return nil
	}
	return nil
}

// Process processes the template with the given content
func (p Processor) Process(tmplContent []byte, values map[string]string) ([]runtime.RawExtension, error) {
	decode := p.codecFactory.UniversalDeserializer().Decode
	obj, _, err := decode(tmplContent, nil, nil)
	if err != nil {
		return nil, errs.Wrapf(err, "unable to decode template")
	}
	tmpl, ok := obj.(*templatev1.Template)
	if !ok {
		return nil, fmt.Errorf("unable to convert object type %T to Template, must be a v1.Template", obj)
	}
	injectUserVars(tmpl, values)
	tmpl, err = processTemplateLocally(tmpl, p.scheme)
	if err != nil {
		return nil, errs.Wrap(err, "unable to process template")
	}
	return tmpl.Objects, nil
}

// injectUserVars injects user specific variables into the Template
func injectUserVars(t *templatev1.Template, values map[string]string) {
	for param, val := range values {
		v := templateprocessing.GetParameterByName(t, param)
		if v != nil {
			v.Value = val
			v.Generate = ""
		}
	}
}

// processTemplateLocally applies the same logic that a remote call would make but makes no
// connection to the server.
func processTemplateLocally(tmpl *templatev1.Template, scheme *runtime.Scheme) (*templatev1.Template, error) {
	processor := templateprocessing.NewProcessor(map[string]generator.Generator{
		"expression": generator.NewExpressionValueGenerator(rand.New(rand.NewSource(time.Now().UnixNano()))),
	})
	if err := processor.Process(tmpl); len(err) > 0 {
		return nil, err.ToAggregate()
	}
	var externalResultObj templatev1.Template
	if err := scheme.Convert(tmpl, &externalResultObj, nil); err != nil {
		return nil, errs.Wrap(err, "failed to convert template to external template object")
	}
	return &externalResultObj, nil
}
