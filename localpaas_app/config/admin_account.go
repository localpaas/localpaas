package config

type AdminAccount struct {
	Email    string `toml:"email" env:"LP_ADMIN_EMAIL"`
	Username string `toml:"username" env:"LP_ADMIN_USERNAME"`
	Password string `toml:"password" env:"LP_ADMIN_PASSWORD"`
}
