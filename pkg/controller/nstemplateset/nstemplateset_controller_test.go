package nstemplateset

import (
	"context"
	"testing"

	"github.com/codeready-toolchain/member-operator/pkg/apis"
	"github.com/codeready-toolchain/toolchain-common/pkg/test"

	"fmt"
	toolchainv1alpha1 "github.com/codeready-toolchain/api/pkg/apis/toolchain/v1alpha1"
	authv1 "github.com/openshift/api/authorization/v1"
	projectv1 "github.com/openshift/api/project/v1"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func TestReconcile(t *testing.T) {

	logf.SetLogger(logf.ZapLogger(true))

	t.Run("create projects", func(t *testing.T) {
		// given
		namespace := uuid.NewV4().String()
		name := uuid.NewV4().String()
		r, req, cl := prepareReconcile(t, namespace, name)
		// also, create the NSTemplateSet CR with the client
		tmplSet := &toolchainv1alpha1.NSTemplateSet{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      name,
			},
			Spec: toolchainv1alpha1.NSTemplateSetSpec{
				TierName: "basic",
			},
		}
		err := cl.Create(context.TODO(), tmplSet)
		require.NoError(t, err)

		// when
		result, err := r.Reconcile(req)
		// then
		require.NoError(t, err)
		assert.Equal(t, reconcile.Result{}, result)
		// check that the project request was created
		verifyProjectRequest(t, r.client, namespace)
		verifyRoleBinding(t, r.client, namespace)
	})

	t.Run("delete_role_binding_and_reconcile", func(t *testing.T) {
		// given
		namespace := uuid.NewV4().String()
		name := uuid.NewV4().String()
		r, req, cl := prepareReconcile(t, namespace, name)
		// also, create the NSTemplateSet CR with the client
		tmplSet := &toolchainv1alpha1.NSTemplateSet{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      name,
			},
			Spec: toolchainv1alpha1.NSTemplateSetSpec{
				TierName: "basic",
			},
		}
		err := cl.Create(context.TODO(), tmplSet)
		require.NoError(t, err)

		// when
		result, err := r.Reconcile(req)

		// then
		require.NoError(t, err)
		assert.Equal(t, reconcile.Result{}, result)

		verifyProjectRequest(t, r.client, namespace)
		verifyRoleBinding(t, r.client, namespace)

		// delete rolebinding to create scenario, of rolebinding failed to create in first run.
		rb, err := getRoleBinding(r.client, namespace)
		require.NoError(t, err)

		err = cl.Delete(context.TODO(), rb)
		require.NoError(t, err)

		result, err = r.Reconcile(req)
		require.NoError(t, err)
		assert.Equal(t, reconcile.Result{}, result)

		verifyProjectRequest(t, r.client, namespace)
		verifyRoleBinding(t, r.client, namespace)
	})

}

func newReconcileRequest(namespace, name string) reconcile.Request {
	return reconcile.Request{
		NamespacedName: types.NamespacedName{
			Namespace: namespace,
			Name:      name,
		},
	}
}

func prepareReconcile(t *testing.T, namespace, name string, initObjs ...runtime.Object) (*ReconcileNSTemplateSet, reconcile.Request, *test.FakeClient) {
	s := scheme.Scheme
	err := apis.AddToScheme(s)
	require.NoError(t, err)
	cl := test.NewFakeClient(t, initObjs...)

	r := &ReconcileNSTemplateSet{
		client: cl,
		scheme: s,
	}
	return r, newReconcileRequest(namespace, name), cl
}

func verifyProjectRequest(t *testing.T, c client.Client, projectRequestName string) {
	// check that the project request was created
	err := c.Get(context.TODO(), types.NamespacedName{Name: projectRequestName, Namespace: ""}, &projectv1.ProjectRequest{}) // project request is cluster-scoped
	require.NoError(t, err)
}

func verifyRoleBinding(t *testing.T, c client.Client, ns string) {
	// check that the rolebinding is created in the namespace
	// (the fake client just records the request but does not perform any consistency check)
	rb, err := getRoleBinding(c, ns)

	require.NotNil(t, rb)
	require.NoError(t, err)
}

func getRoleBinding(c client.Client, ns string) (*authv1.RoleBinding, error) {
	var rb authv1.RoleBinding
	err := c.Get(context.TODO(), types.NamespacedName{
		Namespace: ns,
		Name:      fmt.Sprintf("%s-admin", ns),
	}, &rb)

	return &rb, err
}
