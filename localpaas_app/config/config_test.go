package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_ = os.Setenv("LP_CONFIG_FILE", "testdata/config.myenv.yaml")

		cfg, err := LoadConfig()
		assert.Nil(t, err)
		assert.Equal(t, "myenv", cfg.Env)
		assert.Equal(t, "myplatform", cfg.Platform)
		assert.Equal(t, "Localpaas", cfg.App.Name)
	})

	t.Run("success with override a key with ENV", func(t *testing.T) {
		_ = os.Setenv("LP_CONFIG_FILE", "testdata/config.myenv.yaml")
		_ = os.Setenv("LP_APP_NAME", "overridden")

		cfg, err := LoadConfig()
		assert.Nil(t, err)
		assert.Equal(t, "myenv", cfg.Env)
		assert.Equal(t, "myplatform", cfg.Platform)
		assert.Equal(t, "overridden", cfg.App.Name)
	})

	t.Run("failure: no ENV to find config", func(t *testing.T) {
		_ = os.Unsetenv("LP_ENV")
		_ = os.Unsetenv("LP_CONFIG_FILE")

		_, err := LoadConfig()
		assert.ErrorIs(t, err, ErrConfigFileUnset)
	})

	t.Run("failure: only LP_ENV set, but not found", func(t *testing.T) {
		_ = os.Setenv("LP_ENV", "myenv")
		_ = os.Unsetenv("LP_CONFIG_FILE")

		_, err := LoadConfig()
		assert.ErrorIs(t, err, ErrConfigFileNotFound)
	})

	t.Run("failure: config not found", func(t *testing.T) {
		_ = os.Unsetenv("LP_ENV")
		_ = os.Setenv("LP_CONFIG_FILE", "notexist/config.myenv.yaml")

		_, err := LoadConfig()
		assert.ErrorIs(t, err, ErrConfigFileNotFound)
	})

	t.Run("failure: malformed YAML data", func(t *testing.T) {
		_ = os.Unsetenv("LP_ENV")
		_ = os.Setenv("LP_CONFIG_FILE", "testdata/config-malformed.yaml")

		_, err := LoadConfig()
		assert.ErrorContains(t, err, "yaml: unmarshal errors")
	})
}
