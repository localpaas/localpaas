package entity

import (
	"strings"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/cryptoutil"
)

type SSHKey struct {
	PrivateKey string `json:"privateKey"`
}

func (o *SSHKey) IsEncrypted() bool {
	return strings.HasPrefix(o.PrivateKey, base.SaltPrefix)
}

func (o *SSHKey) Encrypt() error {
	if o.IsEncrypted() {
		return nil
	}
	encrypted, err := cryptoutil.EncryptBase64(o.PrivateKey, base.DefaultSaltLen)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.PrivateKey = encrypted
	return nil
}

func (o *SSHKey) MustEncrypt() *SSHKey {
	gofn.Must1(o.Encrypt())
	return o
}

func (o *SSHKey) Decrypt() error {
	if !o.IsEncrypted() {
		return nil
	}
	decrypted, err := cryptoutil.DecryptBase64(o.PrivateKey)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.PrivateKey = decrypted
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
