package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_ = os.Setenv("LP_CONFIG_FILE", "testdata/config.myenv.toml")

		cfg, err := LoadConfig()
		assert.Nil(t, err)
		assert.Equal(t, "myenv", cfg.Env)
		assert.Equal(t, "myplatform", cfg.Platform)
		assert.Equal(t, "LocalpaaS", cfg.Name)
	})

	t.Run("success with override a key with ENV", func(t *testing.T) {
		_ = os.Setenv("LP_CONFIG_FILE", "testdata/config.myenv.toml")
		_ = os.Setenv("LP_APP_NAME", "overridden")

		cfg, err := LoadConfig()
		assert.Nil(t, err)
		assert.Equal(t, "myenv", cfg.Env)
		assert.Equal(t, "myplatform", cfg.Platform)
		assert.Equal(t, "overridden", cfg.Name)
	})

	t.Run("failure: no ENV to find config", func(t *testing.T) {
		_ = os.Unsetenv("LP_ENV")
		_ = os.Unsetenv("LP_CONFIG_FILE")

		_, err := LoadConfig()
		assert.ErrorIs(t, err, ErrConfigFileUnset)
	})

	t.Run("failure: config not found", func(t *testing.T) {
		_ = os.Unsetenv("LP_ENV")
		_ = os.Setenv("LP_CONFIG_FILE", "notexist/config.myenv.toml")

		_, err := LoadConfig()
		assert.ErrorIs(t, err, ErrConfigFileNotFound)
	})

	t.Run("failure: malformed TOML data", func(t *testing.T) {
		_ = os.Unsetenv("LP_ENV")
		_ = os.Setenv("LP_CONFIG_FILE", "testdata/config-malformed.toml")

		_, err := LoadConfig()
		assert.NotNil(t, err)
	})
}
