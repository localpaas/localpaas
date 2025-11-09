package entity

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
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

	// Salt used to encrypt the secret
	Salt string `json:"salt,omitempty"`
}

func (o *OAuth) IsEncrypted() bool {
	return o.Salt != ""
}

func (o *OAuth) Encrypt() error {
	if o.Salt != "" {
		return nil
	}
	cipher, salt, err := cryptoutil.EncryptBase64(o.ClientSecret, defaultSaltLen)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.ClientSecret = cipher
	o.Salt = salt
	return nil
}

func (o *OAuth) Decrypt() error {
	if o.Salt == "" {
		return nil
	}
	plain, err := cryptoutil.DecryptBase64(o.ClientSecret, o.Salt)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.ClientSecret = plain
	o.Salt = ""
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
