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
	configFileName = "config.toml"
	envPrefix      = "LP"
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

var (
	Current        *Config
	lastConfigFile string
)

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
	SSL        SSL        `toml:"ssl"`
}

func (cfg *Config) IsDevEnv() bool  { return cfg.Env == EnvDev }
func (cfg *Config) IsProdEnv() bool { return cfg.Env == EnvProd }

func (cfg *Config) BaseAPIURL() string {
	return gofn.Must(url.JoinPath(cfg.BaseURL, cfg.HTTPServer.BasePath))
}

/// USER PHOTO

func (cfg *Config) DataPathUserPhoto() string {
	return filepath.Join(cfg.AppPath, "user", "photo")
}
func (cfg *Config) HttpPathUserPhoto() string {
	return "/files/user/photo/"
}

/// SSL CERTS

func (cfg *Config) DataPathCerts() string {
	return filepath.Join(cfg.AppPath, "certs")
}

/// NGINX

func (cfg *Config) DataPathNginx() string {
	return filepath.Join(cfg.AppPath, "nginx")
}
func (cfg *Config) DataPathNginxEtc() string {
	return filepath.Join(cfg.DataPathNginx(), "etc")
}
func (cfg *Config) DataPathNginxEtcConf() string {
	return filepath.Join(cfg.DataPathNginxEtc(), "conf.d")
}
func (cfg *Config) DataPathNginxShare() string {
	return filepath.Join(cfg.DataPathNginx(), "share")
}
func (cfg *Config) DataPathNginxShareDomains() string {
	return filepath.Join(cfg.DataPathNginxShare(), "domains")
}

/// LETS ENCRYPT

func (cfg *Config) DataPathLetsEncrypt() string {
	return filepath.Join(cfg.AppPath, "letsencrypt")
}
func (cfg *Config) DataPathLetsEncryptEtc() string {
	return filepath.Join(cfg.DataPathLetsEncrypt(), "etc")
}

func LoadConfig() (*Config, error) {
	cfg, err := loadConfig("")
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	Current = cfg
	return cfg, nil
}

func loadConfig(configFile string) (*Config, error) {
	config := &Config{}

	if configFile == "" {
		appPath := os.Getenv("LP_APP_PATH")
		if appPath == "" {
			appPath = "/var/lib/localpaas"
		}
		configFile = filepath.Join(appPath, configFileName)

		if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
			configFile = os.Getenv("LP_CONFIG_FILE")
			if configFile == "" {
				return nil, fmt.Errorf("%w: LP_CONFIG_FILE must be defined", ErrConfigFileUnset)
			}
		}
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: %s", ErrConfigFileNotFound, configFile)
	}

	err := configor.New(&configor.Config{ENVPrefix: envPrefix}).Load(config, configFile)
	if err != nil {
		return config, tracerr.Wrap(err)
	}

	lastConfigFile = configFile
	return config, nil
}

func ReloadConfig() (*Config, error) {
	newConfig, err := loadConfig(lastConfigFile)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	// TODO: validate then apply a certain portion of the new config

	Current = newConfig
	return newConfig, nil
}
