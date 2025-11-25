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
}

func (o *Discord) IsEncrypted() bool {
	return strings.HasPrefix(o.Webhook, base.SaltPrefix)
}

func (o *Discord) Encrypt() (*Discord, error) {
	if o.IsEncrypted() {
		return o, nil
	}
	encrypted, err := cryptoutil.EncryptBase64(o.Webhook, base.DefaultSaltLen)
	if err != nil {
		return o, apperrors.Wrap(err)
	}
	o.Webhook = encrypted
	return o, nil
}

func (o *Discord) MustEncrypt() *Discord {
	return gofn.Must(o.Encrypt())
}

func (o *Discord) Decrypt() (*Discord, error) {
	if !o.IsEncrypted() {
		return o, nil
	}
	decrypted, err := cryptoutil.DecryptBase64(o.Webhook)
	if err != nil {
		return o, apperrors.Wrap(err)
	}
	o.Webhook = decrypted
	return o, nil
}

func (o *Discord) MustDecrypt() *Discord {
	return gofn.Must(o.Decrypt())
}

func (s *Setting) AsDiscord() (*Discord, error) {
	if s.parsedData != nil {
		res, ok := s.parsedData.(*Discord)
		if !ok {
			return nil, apperrors.NewTypeInvalid()
		}
		return res, nil
	}
	res := &Discord{}
	if s.Data != "" && s.Type == base.SettingTypeDiscord {
		if err := s.parseData(res); err != nil {
			return nil, apperrors.Wrap(err)
		}
	}
	return res, nil
}

func (s *Setting) MustAsDiscord() *Discord {
	return gofn.Must(s.AsDiscord())
}
