package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

type GithubApp struct {
	ClientID       string         `json:"clientId"`
	ClientSecret   EncryptedField `json:"clientSecret"`
	Organization   string         `json:"org,omitempty"`
	CallbackURL    string         `json:"callbackURL,omitempty"`
	AppID          string         `json:"appId"`
	InstallationID string         `json:"installationId"`
	WebhookURL     string         `json:"webhookURL,omitempty"`
	WebhookSecret  EncryptedField `json:"webhookSecret,omitempty"`
	PrivateKey     EncryptedField `json:"privateKey,omitempty"`
}

func (s *GithubApp) MustDecrypt() *GithubApp {
	s.ClientSecret.MustGetPlain()
	return s
}

func (s *Setting) AsGithubApp() (*GithubApp, error) {
	return parseSettingAs(s, base.SettingTypeGithubApp, func() *GithubApp { return &GithubApp{} })
}

func (s *Setting) MustAsGithubApp() *GithubApp {
	return gofn.Must(s.AsGithubApp())
}
