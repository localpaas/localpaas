package entity

import (
	"encoding/base64"
	"strings"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/cryptoutil"
)

type Secret struct {
	Key    string `json:"k"`
	Value  string `json:"v"`
	Base64 bool   `json:"b64"`
}

func (o *Secret) IsEncrypted() bool {
	return strings.HasPrefix(o.Value, base.SaltPrefix)
}

func (o *Secret) Encrypt() (*Secret, error) {
	if o.IsEncrypted() {
		return o, nil
	}
	encrypted, err := cryptoutil.EncryptBase64(o.Value, base.DefaultSaltLen)
	if err != nil {
		return o, apperrors.Wrap(err)
	}
	o.Value = encrypted
	return o, nil
}

func (o *Secret) MustEncrypt() *Secret {
	return gofn.Must(o.Encrypt())
}

func (o *Secret) Decrypt() (*Secret, error) {
	if !o.IsEncrypted() {
		return o, nil
	}
	decrypted, err := cryptoutil.DecryptBase64(o.Value)
	if err != nil {
		return o, apperrors.Wrap(err)
	}
	o.Value = decrypted
	return o, nil
}

func (o *Secret) MustDecrypt() *Secret {
	return gofn.Must(o.Decrypt())
}

func (o *Secret) ValueAsBytes() []byte {
	o.MustDecrypt()
	if o.Base64 {
		return gofn.Must(base64.StdEncoding.DecodeString(o.Value))
	}
	return []byte(o.Value)
}

func (s *Setting) AsSecret() (*Secret, error) {
	if s.parsedData != nil {
		res, ok := s.parsedData.(*Secret)
		if !ok {
			return nil, apperrors.NewTypeInvalid()
		}
		return res, nil
	}
	res := &Secret{}
	if s.Data != "" && s.Type == base.SettingTypeSecret {
		if err := s.parseData(res); err != nil {
			return nil, apperrors.Wrap(err)
		}
	}
	return res, nil
}

func (s *Setting) MustAsSecret() *Secret {
	return gofn.Must(s.AsSecret())
}
