package template

type MethodHandler struct {
	Methods map[string]MethodDefinition
}

var (
	ResourceMiddlewareStore = map[string]*MethodHandler{
		Namespace:      handlers(handler(Post(After(waitTillProjectActive)))),
		ProjectRequest: handlers(handler(Post(After(waitTillProjectActive)))),
		Project:        handlers(handler(Post(After(waitTillProjectActive)))),
	}
)

func handler(methodDefinition ...func() MethodDefinition) func(methods map[string]MethodDefinition) {
	return func(methods map[string]MethodDefinition) {
		for _, methodDef := range methodDefinition {
			def := methodDef()
			methods[def.method] = def
		}
	}
}

func handlers(actions ...func(methods map[string]MethodDefinition)) *MethodHandler {
	methods := make(map[string]MethodDefinition)
	for _, action := range actions {
		action(methods)
	}
	return &MethodHandler{Methods: methods}
}
