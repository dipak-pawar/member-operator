package template

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func isHandlerPresent(object runtime.Object) (*MethodHandler, bool) {
	kind := object.GetObjectKind().GroupVersionKind().Kind
	handlers, ok := ResourceMiddlewareStore[kind]
	if !ok {
		return nil, false
	}
	return handlers, true
}

func execute(client client.Client, object runtime.Object, method string) error {
	if methodHandler, ok := isHandlerPresent(object); ok {
		return methodHandler.Apply(client, object, method)
	}
	return doExecute(client, object, method)
}

func (e *MethodHandler) Apply(client client.Client, object runtime.Object, method string) error {
	methodDef, ok := e.GetMethodDefinition(method, object)
	// if handler defined on required method
	if ok {
		// execute before handler
		callbackContext := NewHandlerContext(client, object, e, methodDef)
		for _, middleware := range methodDef.beforeHandlers {
			if err := middleware(callbackContext); err != nil {
				return err
			}
		}

		if err := doExecute(client, object, methodDef.method); err != nil {
			return err
		}

		// execute after handler
		for _, middleware := range methodDef.afterHandlers {
			if err := middleware(callbackContext); err != nil {
				return err
			}
		}
	} else { // execute default as handler doesn't exists
		if err := doExecute(client, object, methodDef.method); err != nil {
			return err
		}
	}
	return nil
}

func (e *MethodHandler) GetMethodDefinition(method string, object runtime.Object) (*MethodDefinition, bool) {
	if methodDef, ok := e.Methods[method]; ok {
		return &methodDef, ok
	}
	return nil, false
}

func doExecute(cl client.Client, obj runtime.Object, s string) error {
	//TODO - Check if there is need to return updated object
	switch s {
	case http.MethodPost:
		return cl.Create(context.TODO(), obj)
	case http.MethodDelete:
		return cl.Delete(context.TODO(), obj)
	case http.MethodPatch:
		return cl.Update(context.TODO(), obj)
	}
	return nil
}
