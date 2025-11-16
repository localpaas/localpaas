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
	Passphrase string `json:"passphrase,omitempty"`

	// NOTE: for storing current containing setting only
	Setting *Setting `json:"-"`
}

func (o *SSHKey) IsEncrypted() bool {
	return o.IsPrivateKeyEncrypted() || o.IsPassphraseEncrypted()
}

func (o *SSHKey) Encrypt() error {
	if err := o.EncryptPrivateKey(); err != nil {
		return apperrors.Wrap(err)
	}
	return o.EncryptPassphrase()
}

func (o *SSHKey) MustEncrypt() *SSHKey {
	gofn.Must1(o.Encrypt())
	return o
}

func (o *SSHKey) Decrypt() error {
	if err := o.DecryptPrivateKey(); err != nil {
		return apperrors.Wrap(err)
	}
	return o.DecryptPassphrase()
}

func (o *SSHKey) IsPrivateKeyEncrypted() bool {
	return strings.HasPrefix(o.PrivateKey, base.SaltPrefix)
}

func (o *SSHKey) EncryptPrivateKey() error {
	if o.IsPrivateKeyEncrypted() {
		return nil
	}
	encrypted, err := cryptoutil.EncryptBase64(o.PrivateKey, base.DefaultSaltLen)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.PrivateKey = encrypted
	return nil
}

func (o *SSHKey) MustEncryptPrivateKey() *SSHKey {
	gofn.Must1(o.EncryptPrivateKey())
	return o
}

func (o *SSHKey) DecryptPrivateKey() error {
	if !o.IsPrivateKeyEncrypted() {
		return nil
	}
	decrypted, err := cryptoutil.DecryptBase64(o.PrivateKey)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.PrivateKey = decrypted
	return nil
}

func (o *SSHKey) IsPassphraseEncrypted() bool {
	return strings.HasPrefix(o.Passphrase, base.SaltPrefix)
}

func (o *SSHKey) EncryptPassphrase() error {
	if o.IsPassphraseEncrypted() {
		return nil
	}
	if o.Passphrase == "" {
		return nil
	}
	encrypted, err := cryptoutil.EncryptBase64(o.Passphrase, base.DefaultSaltLen)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.Passphrase = encrypted
	return nil
}

func (o *SSHKey) MustEncryptPassphrase() *SSHKey {
	gofn.Must1(o.EncryptPassphrase())
	return o
}

func (o *SSHKey) DecryptPassphrase() error {
	if !o.IsPassphraseEncrypted() {
		return nil
	}
	if o.Passphrase == "" {
		return nil
	}
	decrypted, err := cryptoutil.DecryptBase64(o.Passphrase)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.Passphrase = decrypted
	return nil
}

func (s *Setting) ParseSSHKey(decrypt bool) (*SSHKey, error) {
	res := &SSHKey{Setting: s}
	if s != nil && s.Data != "" && s.Type == base.SettingTypeSSHKey {
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
	return res, nil
}
