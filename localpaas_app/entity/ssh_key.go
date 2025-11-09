package entity

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/pkg/cryptoutil"
)

type SSHKey struct {
	PrivateKey string `json:"privateKey"`
	Salt       string `json:"salt,omitempty"`
}

func (o *SSHKey) IsEncrypted() bool {
	return o.Salt != ""
}

func (o *SSHKey) Encrypt() error {
	if o.Salt != "" {
		return nil
	}
	cipher, salt, err := cryptoutil.EncryptBase64(o.PrivateKey, defaultSaltLen)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.PrivateKey = cipher
	o.Salt = salt
	return nil
}

func (o *SSHKey) Decrypt() error {
	if o.Salt == "" {
		return nil
	}
	plain, err := cryptoutil.DecryptBase64(o.PrivateKey, o.Salt)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.PrivateKey = plain
	o.Salt = ""
	return nil
}

func (s *Setting) ParseSSHKey(decrypt bool) (*SSHKey, error) {
	if s != nil && s.Data != "" {
		res := &SSHKey{}
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
