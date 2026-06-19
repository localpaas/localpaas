package config

type Agent struct {
	Port        int    `toml:"port" env:"LP_AGENT_PORT" default:"10001"`
	SecretToken string `toml:"secret_token" env:"LP_AGENT_SECRET_TOKEN"`
}
