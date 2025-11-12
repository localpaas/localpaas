package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/jinzhu/configor"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/pkg/tracerr"
)

const (
	configFile = "config.toml"
	envPrefix  = "LP"
)

const (
	PlatformLocal  = "local"
	PlatformRemote = "remote"

	EnvDev  = "development"
	EnvProd = "production"
)

var (
	ErrConfigFileUnset    = errors.New("config file unset")
	ErrConfigFileNotFound = errors.New("config file not found")
)

var Current *Config

type Config struct {
	Env      string `toml:"env" env:"LP_ENV"`
	Platform string `toml:"platform" env:"LP_PLATFORM"`

	Name    string `toml:"name" env:"LP_APP_NAME" default:"LocalPaaS"`
	Version int    `toml:"version" env:"LP_APP_VERSION"`
	BaseURL string `toml:"base_url" env:"LP_APP_BASE_URL"`
	Secret  string `toml:"secret" env:"LP_APP_SECRET" default:"abc123"`
	AppPath string `toml:"app_path" env:"LP_APP_PATH" default:"/var/lib/localpaas"`

	HTTPServer HTTPServer `toml:"http_server"`
	DB         DB         `toml:"db"`
	Cache      Cache      `toml:"cache"`
	Session    Session    `toml:"session"`
}

func (cfg *Config) IsDevEnv() bool  { return cfg.Env == EnvDev }
func (cfg *Config) IsProdEnv() bool { return cfg.Env == EnvProd }

func (cfg *Config) DataPath() string {
	return filepath.Join(cfg.AppPath, "data")
}

func (cfg *Config) DataPathUserPhoto() string {
	return filepath.Join(cfg.DataPath(), "user", "photo")
}

func (cfg *Config) HttpPathUserPhoto() string {
	return "/files/user/photo/"
}

func (cfg *Config) BaseAPIURL() string {
	return gofn.Must(url.JoinPath(cfg.BaseURL, cfg.HTTPServer.BasePath))
}

func LoadConfig() (*Config, error) {
	config := &Config{}

	appPath := os.Getenv("LP_APP_PATH")
	if appPath == "" {
		appPath = "/var/lib/localpaas"
	}
	configFile := filepath.Join(appPath, configFile)

	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		configFile = os.Getenv("LP_CONFIG_FILE")
		if configFile == "" {
			return nil, fmt.Errorf("%w: LP_CONFIG_FILE must be defined", ErrConfigFileUnset)
		}
		if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%w: %s", ErrConfigFileNotFound, configFile)
		}
	}

	cfgr := configor.New(&configor.Config{ENVPrefix: envPrefix})
	if err := cfgr.Load(config, configFile); err != nil {
		return config, tracerr.Wrap(err)
	}

	Current = config
	return config, nil
}
