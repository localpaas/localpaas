package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

type SSHKey struct {
	PrivateKey EncryptedField `json:"privateKey"`
	Passphrase EncryptedField `json:"passphrase,omitzero"`
}

func (o *SSHKey) MustDecrypt() *SSHKey {
	o.PrivateKey.MustGetPlain()
	o.Passphrase.MustGetPlain()
	return o
}

func (s *Setting) AsSSHKey() (*SSHKey, error) {
	return parseSettingAs(s, base.SettingTypeSSHKey, func() *SSHKey { return &SSHKey{} })
}

func (s *Setting) MustAsSSHKey() *SSHKey {
	return gofn.Must(s.AsSSHKey())
}
