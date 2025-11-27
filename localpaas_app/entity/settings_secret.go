package entity

import (
	"encoding/base64"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

type Secret struct {
	Key    string         `json:"k"`
	Value  EncryptedField `json:"v"`
	Base64 bool           `json:"b64"`
}

func (o *Secret) MustDecrypt() *Secret {
	o.Value.MustGetPlain()
	return o
}

func (o *Secret) ValueAsBytes() []byte {
	plain := o.Value.MustGetPlain()
	if o.Base64 {
		return gofn.Must(base64.StdEncoding.DecodeString(plain))
	}
	return []byte(plain)
}

func (s *Setting) AsSecret() (*Secret, error) {
	return parseSettingAs(s, base.SettingTypeSecret, func() *Secret { return &Secret{} })
}

func (s *Setting) MustAsSecret() *Secret {
	return gofn.Must(s.AsSecret())
}
