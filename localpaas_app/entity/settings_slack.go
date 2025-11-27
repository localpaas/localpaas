package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

type Slack struct {
	Webhook EncryptedField `json:"webhook"`
}

func (o *Slack) MustDecrypt() *Slack {
	o.Webhook.MustGetPlain()
	return o
}

func (s *Setting) AsSlack() (*Slack, error) {
	return parseSettingAs(s, base.SettingTypeSlack, func() *Slack { return &Slack{} })
}

func (s *Setting) MustAsSlack() *Slack {
	return gofn.Must(s.AsSlack())
}
