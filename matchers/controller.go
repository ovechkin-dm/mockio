package matchers

import "github.com/ovechkin-dm/mockio/v2/config"


type MockController struct {
	Reporter ErrorReporter
	Config   *config.MockConfig
}

func NewMockController(reporter ErrorReporter, opts ...config.Option) *MockController {
	cfg := config.NewConfig()
	if reporter == nil {
		panic("MockController: provided a nil error reporting (*testing.T) instance")
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return &MockController{
		Reporter: reporter,
		Config:   cfg,
	}	
}
