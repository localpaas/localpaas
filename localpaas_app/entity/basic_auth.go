package entity

import (
	"strings"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/cryptoutil"
)

type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`

	// NOTE: for storing current containing setting only
	Setting *Setting `json:"-"`
}

func (o *BasicAuth) IsEncrypted() bool {
	return strings.HasPrefix(o.Password, base.SaltPrefix)
}

func (o *BasicAuth) Encrypt() error {
	if o.IsEncrypted() {
		return nil
	}
	encrypted, err := cryptoutil.EncryptBase64(o.Password, base.DefaultSaltLen)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.Password = encrypted
	return nil
}

func (o *BasicAuth) MustEncrypt() *BasicAuth {
	gofn.Must1(o.Encrypt())
	return o
}

func (o *BasicAuth) Decrypt() error {
	if !o.IsEncrypted() {
		return nil
	}
	decrypted, err := cryptoutil.DecryptBase64(o.Password)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.Password = decrypted
	return nil
}

func (s *Setting) ParseBasicAuth(decrypt bool) (*BasicAuth, error) {
	res := &BasicAuth{Setting: s}
	if s != nil && s.Data != "" && s.Type == base.SettingTypeBasicAuth {
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
	return res, nil
}
