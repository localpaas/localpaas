package entity

import (
	"strings"

	"github.com/docker/docker/api/types/registry"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/cryptoutil"
)

type RegistryAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

func (o *RegistryAuth) IsEncrypted() bool {
	return strings.HasPrefix(o.Password, base.SaltPrefix)
}

func (o *RegistryAuth) Encrypt() (*RegistryAuth, error) {
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

func (o *RegistryAuth) MustEncrypt() *RegistryAuth {
	return gofn.Must(o.Encrypt())
}

func (o *RegistryAuth) Decrypt() (*RegistryAuth, error) {
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

func (o *RegistryAuth) MustDecrypt() *RegistryAuth {
	return gofn.Must(o.Decrypt())
}

func (o *RegistryAuth) GenerateAuthHeader() (string, error) {
	h, err := registry.EncodeAuthConfig(registry.AuthConfig{
		Username:      o.Username,
		Password:      o.Password,
		ServerAddress: o.Address,
	})
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	return h, nil
}

func (s *Setting) AsRegistryAuth() (*RegistryAuth, error) {
	if s.parsedData != nil {
		res, ok := s.parsedData.(*RegistryAuth)
		if !ok {
			return nil, apperrors.NewTypeInvalid()
		}
		return res, nil
	}
	res := &RegistryAuth{}
	if s.Data != "" && s.Type == base.SettingTypeRegistryAuth {
		if err := s.parseData(res); err != nil {
			return nil, apperrors.Wrap(err)
		}
	}
	return res, nil
}

func (s *Setting) MustAsRegistryAuth() *RegistryAuth {
	return gofn.Must(s.AsRegistryAuth())
}
