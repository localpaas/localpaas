package entity

import (
	"strings"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/cryptoutil"
)

type OAuth struct {
	ClientID     string   `json:"clientId"`
	ClientSecret string   `json:"clientSecret"`
	Organization string   `json:"org,omitempty"`
	CallbackURL  string   `json:"callbackURL,omitempty"`
	AuthURL      string   `json:"authURL,omitempty"`
	TokenURL     string   `json:"tokenURL,omitempty"`
	ProfileURL   string   `json:"profileURL,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
}

func (o *OAuth) IsEncrypted() bool {
	return strings.HasPrefix(o.ClientSecret, base.SaltPrefix)
}

func (o *OAuth) Encrypt() (*OAuth, error) {
	if o.IsEncrypted() {
		return o, nil
	}
	encrypted, err := cryptoutil.EncryptBase64(o.ClientSecret, base.DefaultSaltLen)
	if err != nil {
		return o, apperrors.Wrap(err)
	}
	o.ClientSecret = encrypted
	return o, nil
}

func (o *OAuth) MustEncrypt() *OAuth {
	return gofn.Must(o.Encrypt())
}

func (o *OAuth) Decrypt() (*OAuth, error) {
	if !o.IsEncrypted() {
		return o, nil
	}
	decrypted, err := cryptoutil.DecryptBase64(o.ClientSecret)
	if err != nil {
		return o, apperrors.Wrap(err)
	}
	o.ClientSecret = decrypted
	return o, nil
}

func (o *OAuth) MustDecrypt() *OAuth {
	return gofn.Must(o.Decrypt())
}

func (s *Setting) AsOAuth() (*OAuth, error) {
	if s.parsedData != nil {
		res, ok := s.parsedData.(*OAuth)
		if !ok {
			return nil, apperrors.NewTypeInvalid()
		}
		return res, nil
	}
	res := &OAuth{}
	if s.Data != "" && s.Type == base.SettingTypeOAuth {
		if err := s.parseData(res); err != nil {
			return nil, apperrors.Wrap(err)
		}
	}
	return res, nil
}

func (s *Setting) MustAsOAuth() *OAuth {
	return gofn.Must(s.AsOAuth())
}
