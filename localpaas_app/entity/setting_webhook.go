package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentWebhookVersion = 1
)

type Webhook struct {
	Kind   base.WebhookKind `json:"kind"`
	Secret string           `json:"secret"`
}

func (s *Webhook) GetType() base.SettingType {
	return base.SettingTypeWebhook
}

func (s *Webhook) GetRefSettingIDs() []string {
	return nil
}

func (s *Webhook) MustDecrypt() *Webhook {
	return s
}

func (s *Setting) AsWebhook() (*Webhook, error) {
	return parseSettingAs(s, func() *Webhook { return &Webhook{} })
}

func (s *Setting) MustAsWebhook() *Webhook {
	return gofn.Must(s.AsWebhook())
}
