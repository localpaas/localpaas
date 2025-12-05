package config

import (
	"net/url"
	"path/filepath"

	"github.com/tiendc/gofn"
)

func (cfg *Config) BaseAPIURL() string {
	return gofn.Must(url.JoinPath(cfg.BaseURL, cfg.HTTPServer.BasePath))
}

/// FRONT-END DASHBOARD

func (cfg *Config) DashboardSsoSuccessURL() string {
	return gofn.Must(url.JoinPath(cfg.BaseURL, "auth/sso/success"))
}

/// BACK-END

func (cfg *Config) SsoBaseCallbackURL() string {
	return gofn.Must(url.JoinPath(cfg.BaseAPIURL(), "auth/sso/callback"))
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
