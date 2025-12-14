package entity

import (
	"github.com/docker/docker/api/types/registry"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentRegistryAuthVersion = 1
)

type RegistryAuth struct {
	Username string         `json:"username"`
	Password EncryptedField `json:"password"`
	Address  string         `json:"address"`
}

func (o *RegistryAuth) MustDecrypt() *RegistryAuth {
	o.Password.MustGetPlain()
	return o
}

func (o *RegistryAuth) GenerateAuthHeader() (string, error) {
	password, err := o.Password.GetPlain()
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	h, err := registry.EncodeAuthConfig(registry.AuthConfig{
		Username:      o.Username,
		Password:      password,
		ServerAddress: o.Address,
	})
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	return h, nil
}

func (s *Setting) AsRegistryAuth() (*RegistryAuth, error) {
	return parseSettingAs(s, base.SettingTypeRegistryAuth, func() *RegistryAuth { return &RegistryAuth{} })
}

func (s *Setting) MustAsRegistryAuth() *RegistryAuth {
	return gofn.Must(s.AsRegistryAuth())
}
