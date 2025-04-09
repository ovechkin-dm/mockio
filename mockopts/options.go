package mockopts

import (
	"github.com/ovechkin-dm/mockio/v2/config"
)

// WithoutStackTrace enables stack trace printing for mock errors.
// By default, stack trace is being printed.
// This option is useful for debugging.
// Example:
//
//	SetUp(t, WithoutStackTrace())
func WithoutStackTrace() config.Option {
	return func(cfg *config.MockConfig) {
		cfg.PrintStackTrace = false
	}
}

// StrictVerify enables strict verification of mock calls.
// This means that all mocked methods that are not called will be reported as errors,
// and all not mocked methods that are called will be reported as errors.
// By default, strict verification is disabled.
func StrictVerify() config.Option {
	return func(cfg *config.MockConfig) {
		cfg.StrictVerify = true
	}
}
