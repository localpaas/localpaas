package entity

import (
	"strings"

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

func (o *OAuth) Encrypt() error {
	if o.IsEncrypted() {
		return nil
	}
	encrypted, err := cryptoutil.EncryptBase64(o.ClientSecret, base.DefaultSaltLen)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.ClientSecret = encrypted
	return nil
}

func (o *OAuth) Decrypt() error {
	if !o.IsEncrypted() {
		return nil
	}
	decrypted, err := cryptoutil.DecryptBase64(o.ClientSecret)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.ClientSecret = decrypted
	return nil
}

func (s *Setting) ParseOAuth(decrypt bool) (*OAuth, error) {
	if s != nil && s.Data != "" {
		res := &OAuth{}
		err := s.parseData(res)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		if decrypt {
			if err = res.Decrypt(); err != nil {
				return nil, apperrors.Wrap(err)
			}
		}
		return res, nil
	}
	return nil, nil
}
