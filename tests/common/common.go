package common

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type MockReporter struct {
	reported string
	t        *testing.T
	cleanups []func()
}

func (m *MockReporter) Fatalf(format string, args ...any) {
	m.reported = fmt.Sprintf(format, args...)
}

func (m *MockReporter) IsError() bool {
	return m.reported != ""
}

func (m *MockReporter) ErrorContains(s string) bool {
	return m.IsError() && strings.Contains(strings.ToLower(m.reported), strings.ToLower(s))
}

func (m *MockReporter) GetErrorString() string {
	return m.reported
}

func (m *MockReporter) AssertNoError() {
	if m.IsError() {
		m.t.Fatalf("Expected no error, got: %s", m.reported)
	}
}

func (m *MockReporter) AssertError() {
	if !m.IsError() {
		m.t.Fatalf("Expected error, got nothing")
	}
}

func (m *MockReporter) AssertEqual(expected any, actual any) {
	if !reflect.DeepEqual(expected, actual) {
		m.t.Fatalf("Values not equal. \n Expected: %v \n actual: %v", expected, actual)
	}
}

func (m *MockReporter) AssertErrorContains(err error, s string) {
	if err == nil {
		m.t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), s) {
		m.t.Fatalf("expected error to contain %s", s)
	}
}

func (m *MockReporter) Cleanup(clean func()) {
	m.cleanups = append(m.cleanups, clean)
}

func NewMockReporter(t *testing.T) *MockReporter {
	rep := &MockReporter{
		reported: "",
		t:        t,
		cleanups: make([]func(), 0),
	}
	t.Cleanup(func() {
		for _, v := range rep.cleanups {
			v()
		}
	})
	return rep
}
