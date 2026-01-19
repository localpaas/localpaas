package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentSslVersion = 1
)

type Ssl struct {
	Certificate string           `json:"certificate"`
	PrivateKey  EncryptedField   `json:"privateKey"`
	KeySize     int              `json:"keySize"`
	Provider    base.SslProvider `json:"provider,omitempty"`
	Email       string           `json:"email"`
}

func (o *Ssl) MustDecrypt() *Ssl {
	o.PrivateKey.MustGetPlain()
	return o
}

func (s *Setting) AsSsl() (*Ssl, error) {
	return parseSettingAs(s, base.SettingTypeSSL, func() *Ssl { return &Ssl{} })
}

func (s *Setting) MustAsSsl() *Ssl {
	return gofn.Must(s.AsSsl())
}
