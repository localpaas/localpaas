package entity

import (
	"strings"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/cryptoutil"
)

type Discord struct {
	Webhook string `json:"webhook"`

	// NOTE: for storing current containing setting only
	Setting *Setting `json:"-"`
}

func (o *Discord) IsEncrypted() bool {
	return strings.HasPrefix(o.Webhook, base.SaltPrefix)
}

func (o *Discord) Encrypt() error {
	if o.IsEncrypted() {
		return nil
	}
	encrypted, err := cryptoutil.EncryptBase64(o.Webhook, base.DefaultSaltLen)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.Webhook = encrypted
	return nil
}

func (o *Discord) MustEncrypt() *Discord {
	gofn.Must1(o.Encrypt())
	return o
}

func (o *Discord) Decrypt() error {
	if !o.IsEncrypted() {
		return nil
	}
	decrypted, err := cryptoutil.DecryptBase64(o.Webhook)
	if err != nil {
		return apperrors.Wrap(err)
	}
	o.Webhook = decrypted
	return nil
}

func (s *Setting) ParseDiscord(decrypt bool) (*Discord, error) {
	res := &Discord{Setting: s}
	if s != nil && s.Data != "" && s.Type == base.SettingTypeDiscord {
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
