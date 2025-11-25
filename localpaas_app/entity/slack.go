package entity

import (
	"strings"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/cryptoutil"
)

type Slack struct {
	Webhook string `json:"webhook"`
}

func (o *Slack) IsEncrypted() bool {
	return strings.HasPrefix(o.Webhook, base.SaltPrefix)
}

func (o *Slack) Encrypt() (*Slack, error) {
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

func (o *Slack) MustEncrypt() *Slack {
	return gofn.Must(o.Encrypt())
}

func (o *Slack) Decrypt() (*Slack, error) {
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

func (o *Slack) MustDecrypt() *Slack {
	return gofn.Must(o.Decrypt())
}

func (s *Setting) AsSlack() (*Slack, error) {
	if s.parsedData != nil {
		res, ok := s.parsedData.(*Slack)
		if !ok {
			return nil, apperrors.NewTypeInvalid()
		}
		return res, nil
	}
	res := &Slack{}
	if s.Data != "" && s.Type == base.SettingTypeSlack {
		if err := s.parseData(res); err != nil {
			return nil, apperrors.Wrap(err)
		}
	}
	return res, nil
}

func (s *Setting) MustAsSlack() *Slack {
	return gofn.Must(s.AsSlack())
}
