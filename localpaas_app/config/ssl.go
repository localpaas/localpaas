package config

type SSL struct {
	LeUserEmail string `toml:"le_user_email" env:"LP_SSL_LE_USER_EMAIL"`
}
