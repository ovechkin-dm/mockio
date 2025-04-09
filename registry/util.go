package registry

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
)

const (
	PackageName     = "github.com/ovechkin-dm/mockio"
	DynoPackageName = "github.com/ovechkin-dm/go-dyno"
	TestPackageName = "github.com/ovechkin-dm/mockio/v2/tests"
	DebugPackage    = "runtime/debug.Stack()"
	GOIDPackageName = "github.com/petermattis/goid"
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
		result[i] = valueToInterface(values[i])
	}
	return result
}

func valueToInterface(value reflect.Value) any {
	return value.Interface()
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

type StackTrace struct {
	goroutine string
	lines     []*StackLine
}

type StackLine struct {
	Path string
	Line string
}

func (s *StackLine) String() string {
	return fmt.Sprintf("%s\n\t%s", s.Path, s.Line)
}

func (s *StackLine) IsLibraryStackLine() bool {
	if strings.Contains(s.Path, DebugPackage) {
		return true
	}
	if strings.Contains(s.Path, DynoPackageName) {
		return true
	}
	if strings.Contains(s.Path, GOIDPackageName) {
		return true
	}
	return strings.Contains(s.Path, PackageName) && !strings.Contains(s.Path, TestPackageName)
}

func (s *StackTrace) String() string {
	result := make([]string, 0)
	for i := range s.lines {
		result = append(result, s.lines[i].String()+"\n")
	}
	return strings.Join(result, "")
}

func (s *StackTrace) CallerLine() string {
	for i := range s.lines {
		if s.lines[i].IsLibraryStackLine() {
			if i < len(s.lines)-1 && !s.lines[i+1].IsLibraryStackLine() {
				return s.lines[i+1].Line
			}
		}
	}
	return ""
}

func (s *StackTrace) WithoutLibraryCalls() *StackTrace {
	var result []*StackLine
	for i := range s.lines {
		if !s.lines[i].IsLibraryStackLine() {
			result = append(result, s.lines[i])
		}
	}
	return &StackTrace{
		lines: result,
	}
}

func NewStackTrace() *StackTrace {
	stack := string(debug.Stack())
	lines := strings.Split(stack, "\n")
	goroutine := lines[0]
	stackLines := make([]*StackLine, 0)
	for i := 1; i < len(lines)-1; i += 2 {
		l := &StackLine{
			Path: strings.TrimSpace(lines[i]),
			Line: strings.TrimSpace(lines[i+1]),
		}
		stackLines = append(stackLines, l)
	}
	return &StackTrace{
		goroutine: goroutine,
		lines:     stackLines,
	}
}
