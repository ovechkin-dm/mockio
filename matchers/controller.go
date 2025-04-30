package matchers

import (
	"reflect"

	"github.com/ovechkin-dm/mockio/v2/config"
)

type MockEnv struct {
	Reporter ErrorReporter
	Config   *config.MockConfig
}

type MockController struct {
	Env         *MockEnv
	MockFactory MockFactory
}

type MockFactory interface {
	BuildHandler(env *MockEnv, ifaceType reflect.Type) Handler
}

type Handler interface {
	Handle(method reflect.Method, values []reflect.Value) []reflect.Value
}
