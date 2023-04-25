package registry

import (
	"reflect"
)

func createDefaultReturnValues(m reflect.Method) []reflect.Value {
	result := make([]reflect.Value, m.Type.NumOut())
	for i := 0; i < m.Type.NumOut(); i++ {
		result[i] = reflect.New(m.Type.Out(i)).Elem()
	}
	return result
}

func valueSliceToInterfaceSlice(values []reflect.Value) []any {
	result := make([]any, len(values))
	for i := range values {
		result[i] = values[i].Interface()
	}
	return result
}

func interfaceSliceToValueSlice(values []any, m reflect.Method) []reflect.Value {
	result := make([]reflect.Value, len(values))
	for i := range values {
		retV := reflect.New(m.Type.Out(i)).Elem()
		if values[i] != nil {
			retV.Set(reflect.ValueOf(values[i]))
		}
		result[i] = retV
	}
	return result
}
