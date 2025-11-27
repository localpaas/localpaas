package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

type OAuth struct {
	ClientID     string         `json:"clientId"`
	ClientSecret EncryptedField `json:"clientSecret"`
	Organization string         `json:"org,omitempty"`
	CallbackURL  string         `json:"callbackURL,omitempty"`
	AuthURL      string         `json:"authURL,omitempty"`
	TokenURL     string         `json:"tokenURL,omitempty"`
	ProfileURL   string         `json:"profileURL,omitempty"`
	Scopes       []string       `json:"scopes,omitempty"`
}

func (s *OAuth) MustDecrypt() *OAuth {
	s.ClientSecret.MustGetPlain()
	return s
}

func (s *Setting) AsOAuth() (*OAuth, error) {
	return parseSettingAs(s, base.SettingTypeOAuth, func() *OAuth { return &OAuth{} })
}

func (s *Setting) MustAsOAuth() *OAuth {
	return gofn.Must(s.AsOAuth())
}
