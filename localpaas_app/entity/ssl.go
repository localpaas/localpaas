package entity

import (
	"strings"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/cryptoutil"
)

type Ssl struct {
	Certificate string           `json:"certificate"`
	PrivateKey  string           `json:"privateKey"`
	KeySize     int              `json:"keySize"`
	Provider    base.SslProvider `json:"provider,omitempty"`
	Email       string           `json:"email"`
	Expiration  time.Time        `json:"expiration,omitzero"`

	// NOTE: for storing current containing setting only
	Setting *Setting `json:"-"`
}

func (o *Ssl) IsEncrypted() bool {
	return strings.HasPrefix(o.PrivateKey, base.SaltPrefix)
}

func (o *Ssl) Encrypt() error {
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

func (o *Ssl) MustEncrypt() *Ssl {
	gofn.Must1(o.Encrypt())
	return o
}

func (o *Ssl) Decrypt() error {
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

func (s *Setting) ParseSsl(decrypt bool) (*Ssl, error) {
	res := &Ssl{Setting: s}
	if s != nil && s.Data != "" && s.Type == base.SettingTypeSsl {
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
