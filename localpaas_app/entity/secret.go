package entity

import (
	"strings"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/cryptoutil"
)

type Secret struct {
	Key   string `json:"k"`
	Value string `json:"v"`
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

func (s *Setting) ParseSecret(decrypt bool) (*Secret, error) {
	if s != nil && s.Data != "" {
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
