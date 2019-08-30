package template

import (
	"net/http"
)

type MethodDefinition struct {
	method         string
	beforeHandlers []Handler
	afterHandlers  []Handler
}

func NewMethodDefinition(action string, beforeHandlers []Handler, afterHandlers []Handler, modifiers ...MethodDefModifier) MethodDefinition {
	methodDefinition := MethodDefinition{
		method:         action,
		beforeHandlers: beforeHandlers,
		afterHandlers:  afterHandlers,
	}
	for _, modify := range modifiers {
		modify(&methodDefinition)
	}
	return methodDefinition
}

func newMethodDefinition(action string, modifiers ...MethodDefModifier) MethodDefinition {
	var before, after []Handler
	return NewMethodDefinition(action, before, after, modifiers...)
}

type MethodDefModifier func(*MethodDefinition)

func Post(modifiers ...MethodDefModifier) func() MethodDefinition {
	return func() MethodDefinition {
		return newMethodDefinition(http.MethodPost, modifiers...)
	}
}

func Get(modifiers ...MethodDefModifier) func() MethodDefinition {
	return func() MethodDefinition {
		return newMethodDefinition(http.MethodGet, modifiers...)
	}
}

func Delete(modifiers ...MethodDefModifier) func() MethodDefinition {
	return func() MethodDefinition {
		return newMethodDefinition(http.MethodDelete, modifiers...)
	}
}

func Patch(modifiers ...MethodDefModifier) func() MethodDefinition {
	return func() MethodDefinition {
		return newMethodDefinition(http.MethodPatch, modifiers...)
	}
}
