package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jinzhu/configor"

	"github.com/localpaas/localpaas/pkg/tracerr"
)

const (
	configFileTemplate = "config/config.{{env}}.yaml"
	envPrefix          = "LP"
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
	Env      string  `yaml:"env" env:"LP_ENV"`
	Platform string  `yaml:"platform" env:"LP_PLATFORM"`
	DevMode  DevMode `yaml:"dev_mode"`

	App        App        `yaml:"app"`
	HTTPServer HTTPServer `yaml:"http_server"`
	DB         DB         `yaml:"db"`
	Cache      Cache      `yaml:"cache"`
	Session    Session    `yaml:"session"`
}

func (cfg *Config) IsDevEnv() bool  { return cfg.Env == EnvDev }
func (cfg *Config) IsProdEnv() bool { return cfg.Env == EnvProd }

func LoadConfig() (*Config, error) {
	config := &Config{}

	// Finds config file based on 2 ENVs: LP_ENV and LP_CONFIG_FILE
	// with LP_CONFIG_FILE has higher priority
	configFile := os.Getenv("LP_CONFIG_FILE")
	if configFile == "" {
		env := os.Getenv("LP_ENV")
		if env == "" {
			return nil, fmt.Errorf("%w: either LP_ENV or LP_CONFIG_FILE must be defined", ErrConfigFileUnset)
		}

		platform := os.Getenv("LP_PLATFORM")
		if platform == "" {
			configFile = strings.Replace(configFileTemplate, "{{env}}", PlatformLocal, 1)
		} else {
			configFile = strings.Replace(configFileTemplate, "{{env}}", strings.ToLower(env), 1)
		}
	}

	// configor doesn't check file existence
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("%w: %s", ErrConfigFileNotFound, configFile)
	}

	cfgr := configor.New(&configor.Config{ENVPrefix: envPrefix})
	if err := cfgr.Load(config, configFile); err != nil {
		return config, tracerr.Wrap(err)
	}

	Current = config
	return config, nil
}
