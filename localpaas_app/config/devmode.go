package config

type DevMode struct {
	SkipAuthCheck bool `yaml:"skip_auth_check" env:"LP_DEV_MODE_SKIP_AUTH_CHECK"`
}
