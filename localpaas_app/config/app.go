package config

type App struct {
	Name    string `yaml:"name" env:"LP_APP_NAME" default:"LocalPaaS"`
	Version int    `yaml:"version" env:"LP_APP_VERSION"`
	BaseURL string `yaml:"base_url" env:"LP_APP_BASE_URL"`
	Secret  string `yaml:"secret" env:"LP_APP_SECRET" default:"abc123"`
}
