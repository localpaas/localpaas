package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentSSHKeyVersion = 1
)

type SSHKey struct {
	PrivateKey EncryptedField `json:"privateKey"`
	Passphrase EncryptedField `json:"passphrase,omitzero"`
}

func (s *SSHKey) GetType() base.SettingType {
	return base.SettingTypeSSHKey
}

func (s *SSHKey) GetRefSettingIDs() []string {
	return nil
}

func (s *SSHKey) MustDecrypt() *SSHKey {
	s.PrivateKey.MustGetPlain()
	s.Passphrase.MustGetPlain()
	return s
}

func (s *Setting) AsSSHKey() (*SSHKey, error) {
	return parseSettingAs(s, func() *SSHKey { return &SSHKey{} })
}

func (s *Setting) MustAsSSHKey() *SSHKey {
	return gofn.Must(s.AsSSHKey())
}
