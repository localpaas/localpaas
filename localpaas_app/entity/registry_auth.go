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

	// NOTE: for storing current containing setting only
	Setting *Setting `json:"-"`
}

func (o *RegistryAuth) IsEncrypted() bool {
	return strings.HasPrefix(o.Password, base.SaltPrefix)
}

func (o *RegistryAuth) Encrypt() error {
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

func (o *RegistryAuth) MustEncrypt() *RegistryAuth {
	gofn.Must1(o.Encrypt())
	return o
}

func (o *RegistryAuth) Decrypt() error {
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

func (s *Setting) ParseRegistryAuth(decrypt bool) (*RegistryAuth, error) {
	res := &RegistryAuth{Setting: s}
	if s != nil && s.Data != "" && s.Type == base.SettingTypeRegistryAuth {
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
