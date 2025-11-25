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
}

func (o *SSHKey) IsEncrypted() bool {
	return o.IsPrivateKeyEncrypted() || o.IsPassphraseEncrypted()
}

func (o *SSHKey) Encrypt() (*SSHKey, error) {
	if _, err := o.EncryptPrivateKey(); err != nil {
		return o, apperrors.Wrap(err)
	}
	return o.EncryptPassphrase()
}

func (o *SSHKey) MustEncrypt() *SSHKey {
	return gofn.Must(o.Encrypt())
}

func (o *SSHKey) Decrypt() (*SSHKey, error) {
	if _, err := o.DecryptPrivateKey(); err != nil {
		return o, apperrors.Wrap(err)
	}
	return o.DecryptPassphrase()
}

func (o *SSHKey) IsPrivateKeyEncrypted() bool {
	return strings.HasPrefix(o.PrivateKey, base.SaltPrefix)
}

func (o *SSHKey) EncryptPrivateKey() (*SSHKey, error) {
	if o.IsPrivateKeyEncrypted() {
		return o, nil
	}
	encrypted, err := cryptoutil.EncryptBase64(o.PrivateKey, base.DefaultSaltLen)
	if err != nil {
		return o, apperrors.Wrap(err)
	}
	o.PrivateKey = encrypted
	return o, nil
}

func (o *SSHKey) MustEncryptPrivateKey() *SSHKey {
	return gofn.Must(o.EncryptPrivateKey())
}

func (o *SSHKey) DecryptPrivateKey() (*SSHKey, error) {
	if !o.IsPrivateKeyEncrypted() {
		return o, nil
	}
	decrypted, err := cryptoutil.DecryptBase64(o.PrivateKey)
	if err != nil {
		return o, apperrors.Wrap(err)
	}
	o.PrivateKey = decrypted
	return o, nil
}

func (o *SSHKey) IsPassphraseEncrypted() bool {
	return strings.HasPrefix(o.Passphrase, base.SaltPrefix)
}

func (o *SSHKey) EncryptPassphrase() (*SSHKey, error) {
	if o.IsPassphraseEncrypted() {
		return o, nil
	}
	if o.Passphrase == "" {
		return o, nil
	}
	encrypted, err := cryptoutil.EncryptBase64(o.Passphrase, base.DefaultSaltLen)
	if err != nil {
		return o, apperrors.Wrap(err)
	}
	o.Passphrase = encrypted
	return o, nil
}

func (o *SSHKey) MustEncryptPassphrase() *SSHKey {
	return gofn.Must(o.EncryptPassphrase())
}

func (o *SSHKey) DecryptPassphrase() (*SSHKey, error) {
	if !o.IsPassphraseEncrypted() {
		return o, nil
	}
	if o.Passphrase == "" {
		return o, nil
	}
	decrypted, err := cryptoutil.DecryptBase64(o.Passphrase)
	if err != nil {
		return o, apperrors.Wrap(err)
	}
	o.Passphrase = decrypted
	return o, nil
}

func (o *SSHKey) MustDecrypt() *SSHKey {
	return gofn.Must(o.Decrypt())
}

func (s *Setting) AsSSHKey() (*SSHKey, error) {
	if s.parsedData != nil {
		res, ok := s.parsedData.(*SSHKey)
		if !ok {
			return nil, apperrors.NewTypeInvalid()
		}
		return res, nil
	}
	res := &SSHKey{}
	if s.Data != "" && s.Type == base.SettingTypeSSHKey {
		if err := s.parseData(res); err != nil {
			return nil, apperrors.Wrap(err)
		}
	}
	return res, nil
}

func (s *Setting) MustAsSSHKey() *SSHKey {
	return gofn.Must(s.AsSSHKey())
}
