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
}

func (o *Ssl) IsEncrypted() bool {
	return strings.HasPrefix(o.PrivateKey, base.SaltPrefix)
}

func (o *Ssl) Encrypt() (*Ssl, error) {
	if o.IsEncrypted() {
		return o, nil
	}
	encrypted, err := cryptoutil.EncryptBase64(o.PrivateKey, base.DefaultSaltLen)
	if err != nil {
		return o, apperrors.Wrap(err)
	}
	o.PrivateKey = encrypted
	return o, nil
}

func (o *Ssl) MustEncrypt() *Ssl {
	return gofn.Must(o.Encrypt())
}

func (o *Ssl) Decrypt() (*Ssl, error) {
	if !o.IsEncrypted() {
		return o, nil
	}
	decrypted, err := cryptoutil.DecryptBase64(o.PrivateKey)
	if err != nil {
		return o, apperrors.Wrap(err)
	}
	o.PrivateKey = decrypted
	return o, nil
}

func (o *Ssl) MustDecrypt() *Ssl {
	return gofn.Must(o.Decrypt())
}

func (s *Setting) AsSsl() (*Ssl, error) {
	if s.parsedData != nil {
		res, ok := s.parsedData.(*Ssl)
		if !ok {
			return nil, apperrors.NewTypeInvalid()
		}
		return res, nil
	}
	res := &Ssl{}
	if s.Data != "" && s.Type == base.SettingTypeSsl {
		if err := s.parseData(res); err != nil {
			return nil, apperrors.Wrap(err)
		}
	}
	return res, nil
}

func (s *Setting) MustAsSsl() *Ssl {
	return gofn.Must(s.AsSsl())
}
