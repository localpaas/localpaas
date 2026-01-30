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

func (s *RegistryAuth) GetType() base.SettingType {
	return base.SettingTypeRegistryAuth
}

func (s *RegistryAuth) GetRefSettingIDs() []string {
	return nil
}

func (s *RegistryAuth) MustDecrypt() *RegistryAuth {
	s.Password.MustGetPlain()
	return s
}

func (s *RegistryAuth) GenerateAuthHeader() (string, error) {
	password, err := s.Password.GetPlain()
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	h, err := registry.EncodeAuthConfig(registry.AuthConfig{
		Username:      s.Username,
		Password:      password,
		ServerAddress: s.Address,
	})
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	return h, nil
}

func (s *Setting) AsRegistryAuth() (*RegistryAuth, error) {
	return parseSettingAs(s, func() *RegistryAuth { return &RegistryAuth{} })
}

func (s *Setting) MustAsRegistryAuth() *RegistryAuth {
	return gofn.Must(s.AsRegistryAuth())
}
