package template

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	v1 "github.com/openshift/api/project/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type Handler func(ctx HandlerContext) error

type HandlerContext struct {
	client        client.Client
	object        runtime.Object
	methodHandler *MethodHandler
	method        *MethodDefinition
}

func NewHandlerContext(client client.Client, object runtime.Object, objEndpoints *MethodHandler, method *MethodDefinition) HandlerContext {
	return HandlerContext{
		client:        client,
		object:        object,
		methodHandler: objEndpoints,
		method:        method,
	}
}

func Before(handlers ...Handler) MethodDefModifier {
	return func(methodDefinition *MethodDefinition) {
		methodDefinition.beforeHandlers = append(methodDefinition.beforeHandlers, handlers...)
	}
}

func After(handlers ...Handler) MethodDefModifier {
	return func(methodDefinition *MethodDefinition){
		methodDefinition.afterHandlers = append(methodDefinition.afterHandlers, handlers...)
	}
}

var waitTillProjectActive = func(ctx HandlerContext) error {
	objKey, err := client.ObjectKeyFromObject(ctx.object)
	if err != nil {
		return err
	}
	return wait.Poll(1*time.Second, 30*time.Second, func() (done bool, err error) {
		var prj v1.Project
		fmt.Println("inside poll")
		if err = ctx.client.Get(context.Background(), objKey, &prj); err != nil {
			return false, nil
		}

		fmt.Println("prj.status.phase", prj.Status.Phase)
		if prj.Status.Phase == corev1.NamespaceActive {
			return true, nil
		}
		return false, nil
	})
}
