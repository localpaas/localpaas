package config

import (
	"fmt"
	"net/url"

	"github.com/tiendc/gofn"
)

type HTTPServer struct {
	BaseURL  string `yaml:"base_url" env:"LP_HTTP_SERVER_BASE_URL"`
	BasePath string `yaml:"base_path" env:"LP_HTTP_SERVER_BASE_PATH" default:"/api/v1"`
	Port     int    `yaml:"port" env:"LP_HTTP_SERVER_PORT"`
	CORS     struct {
		AllowOrigins []string `yaml:"allow_origins" env:"LP_HTTP_SERVER_CORS_ALLOW_ORIGINS"`
	} `yaml:"cors"`
}

func (c *HTTPServer) BindingAddress() string {
	return fmt.Sprintf(":%d", c.Port)
}

func (c *HTTPServer) BaseAPIURL() string {
	return gofn.Must(url.JoinPath(c.BaseURL, c.BasePath))
}
