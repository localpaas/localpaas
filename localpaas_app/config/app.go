package config

import "path/filepath"

type App struct {
	Name    string `yaml:"name" env:"LP_APP_NAME" default:"LocalPaaS"`
	Version int    `yaml:"version" env:"LP_APP_VERSION"`
	BaseURL string `yaml:"base_url" env:"LP_APP_BASE_URL"`
	Secret  string `yaml:"secret" env:"LP_APP_SECRET" default:"abc123"`
	AppPath string `yaml:"app_path" env:"LP_APP_PATH" default:"/var/lib/localpaas"`
}

func (a *App) DataPath() string {
	return filepath.Join(a.AppPath, "data")
}

func (a *App) DataPathUserPhoto() string {
	return filepath.Join(a.DataPath(), "user", "photo")
}

func (a *App) HttpPathUserPhoto() string {
	return "/files/user/photo/"
}
