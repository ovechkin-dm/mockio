package config

type Option func(*MockConfig)

type MockConfig struct {
	PrintStackTrace bool
	StrictVerify    bool
}

func NewConfig() *MockConfig {
	return &MockConfig{
		PrintStackTrace: true,
		StrictVerify:    false,
	}
}
