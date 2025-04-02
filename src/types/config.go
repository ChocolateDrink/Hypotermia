package types

type ConfigConfig struct {
	Debugging DebuggingConfig `json:"debugging"`
}

type DebuggingConfig struct {
	DebuggingEnabled bool `json:"debugging_enabled"`
	Persistence      bool `json:"persistence"`
	Hidden           bool `json:"hidden"`
}
