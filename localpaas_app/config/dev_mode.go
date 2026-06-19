package config

type DevMode struct {
	Enabled         bool `toml:"-" env:"-"`
	ForceAgentLocal bool `toml:"force_agent_local" env:"LP_DEV_MODE_FORCE_AGENT_LOCAL"`
}
