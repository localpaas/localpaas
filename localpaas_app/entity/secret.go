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

func (o *Secret) Encrypt() error {
	if o.IsEncrypted() {
		return nil
	}
	encrypted, err := cryptoutil.EncryptBase64(o.Value, base.DefaultSaltLen)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.Value = encrypted
	return nil
}

func (o *Secret) MustEncrypt() *Secret {
	gofn.Must1(o.Encrypt())
	return o
}

func (o *Secret) Decrypt() error {
	if !o.IsEncrypted() {
		return nil
	}
	decrypted, err := cryptoutil.DecryptBase64(o.Value)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.Value = decrypted
	return nil
}

func (o *Secret) ValueAsBytes() []byte {
	gofn.Must1(o.Decrypt())
	if o.Base64 {
		return gofn.Must(base64.StdEncoding.DecodeString(o.Value))
	}
	return []byte(o.Value)
}

func (s *Setting) ParseSecret(decrypt bool) (*Secret, error) {
	if s != nil && s.Data != "" && s.Type == base.SettingTypeSecret {
		res := &Secret{}
		err := s.parseData(res)
		if err != nil {
			return nil, err
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
