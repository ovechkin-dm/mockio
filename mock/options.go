package mock

import (
	"github.com/ovechkin-dm/mockio/config"
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
