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
}

func (o *BasicAuth) IsEncrypted() bool {
	return strings.HasPrefix(o.Password, base.SaltPrefix)
}

func (o *BasicAuth) Encrypt() (*BasicAuth, error) {
	if o.IsEncrypted() {
		return o, nil
	}
	encrypted, err := cryptoutil.EncryptBase64(o.Password, base.DefaultSaltLen)
	if err != nil {
		return o, apperrors.Wrap(err)
	}
	o.Password = encrypted
	return o, nil
}

func (o *BasicAuth) MustEncrypt() *BasicAuth {
	return gofn.Must(o.Encrypt())
}

func (o *BasicAuth) Decrypt() (*BasicAuth, error) {
	if !o.IsEncrypted() {
		return o, nil
	}
	decrypted, err := cryptoutil.DecryptBase64(o.Password)
	if err != nil {
		return o, apperrors.Wrap(err)
	}
	o.Password = decrypted
	return o, nil
}

func (o *BasicAuth) MustDecrypt() *BasicAuth {
	return gofn.Must(o.Decrypt())
}

func (s *Setting) AsBasicAuth() (*BasicAuth, error) {
	if s.parsedData != nil {
		res, ok := s.parsedData.(*BasicAuth)
		if !ok {
			return nil, apperrors.NewTypeInvalid()
		}
		return res, nil
	}
	res := &BasicAuth{}
	if s.Data != "" && s.Type == base.SettingTypeBasicAuth {
		if err := s.parseData(res); err != nil {
			return nil, apperrors.Wrap(err)
		}
	}
	return res, nil
}

func (s *Setting) MustAsBasicAuth() *BasicAuth {
	return gofn.Must(s.AsBasicAuth())
}
