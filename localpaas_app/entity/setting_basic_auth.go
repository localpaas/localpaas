package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentBasicAuthVersion = 1
)

type BasicAuth struct {
	Username string         `json:"username"`
	Password EncryptedField `json:"password"`
}

func (s *BasicAuth) MustDecrypt() *BasicAuth {
	s.Password.MustGetPlain()
	return s
}

func (s *Setting) AsBasicAuth() (*BasicAuth, error) {
	return parseSettingAs(s, base.SettingTypeBasicAuth, func() *BasicAuth { return &BasicAuth{} })
}

func (s *Setting) MustAsBasicAuth() *BasicAuth {
	return gofn.Must(s.AsBasicAuth())
}
