package config

import (
	"fmt"
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

func (cfg *Config) DashboardUserSignupURL(token string) string {
	return gofn.Must(url.JoinPath(cfg.BaseURL, "auth/sign-up")) +
		fmt.Sprintf("?token=%s", token)
}

func (cfg *Config) DashboardPasswordResetURL(userID, token string) string {
	return gofn.Must(url.JoinPath(cfg.BaseURL, "auth/reset-password")) +
		fmt.Sprintf("?userId=%s&token=%s", userID, token)
}

func (cfg *Config) DashboardDeploymentDetailsURL(deploymentID string) string {
	return gofn.Must(url.JoinPath(cfg.BaseURL, "deployments", deploymentID)) // TODO: update this
}

func (cfg *Config) DashboardCronTaskDetailsURL(cronJobID, taskID string) string {
	return gofn.Must(url.JoinPath(cfg.BaseURL, "cron-jobs", cronJobID, "tasks", taskID)) // TODO: update this
}

func (cfg *Config) DashboardHealthcheckDetailsURL(settingID, taskID string) string {
	return gofn.Must(url.JoinPath(cfg.BaseURL, "healthcheck", settingID, "tasks", taskID)) // TODO: update this
}

/// BACK-END

func (cfg *Config) SsoBaseCallbackURL() string {
	return gofn.Must(url.JoinPath(cfg.BaseAPIURL(), "auth/sso/callback"))
}

/// OBJECT PHOTOS

func (cfg *Config) DataPathPhoto() string {
	return filepath.Join(cfg.AppPath, "files", "photo")
}
func (cfg *Config) HttpPathPhoto() string {
	return "/files/photo/"
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
func (cfg *Config) DataPathNginxShareBasicAuth() string {
	return filepath.Join(cfg.DataPathNginxShare(), "basic-auth")
}

/// LETS ENCRYPT

func (cfg *Config) DataPathLetsEncrypt() string {
	return filepath.Join(cfg.AppPath, "letsencrypt")
}
func (cfg *Config) DataPathLetsEncryptEtc() string {
	return filepath.Join(cfg.DataPathLetsEncrypt(), "etc")
}
