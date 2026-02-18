package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentSSHKeyVersion = 1
)

var _ = registerSettingParser(base.SettingTypeSSHKey, &sshKeyParser{})

type sshKeyParser struct {
}

func (s *sshKeyParser) New() SettingData {
	return &SSHKey{}
}

type SSHKey struct {
	PrivateKey EncryptedField `json:"privateKey"`
	Passphrase EncryptedField `json:"passphrase,omitzero"`
}

func (s *SSHKey) GetType() base.SettingType {
	return base.SettingTypeSSHKey
}

func (s *SSHKey) GetRefObjectIDs() *RefObjectIDs {
	return &RefObjectIDs{}
}

func (s *SSHKey) MustDecrypt() *SSHKey {
	s.PrivateKey.MustGetPlain()
	s.Passphrase.MustGetPlain()
	return s
}

func (s *Setting) AsSSHKey() (*SSHKey, error) {
	return parseSettingAs[*SSHKey](s)
}

func (s *Setting) MustAsSSHKey() *SSHKey {
	return gofn.Must(s.AsSSHKey())
}
