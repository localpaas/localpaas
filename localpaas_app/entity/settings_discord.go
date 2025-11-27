package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

type Discord struct {
	Webhook EncryptedField `json:"webhook"`
}

func (s *Discord) MustDecrypt() *Discord {
	s.Webhook.MustGetPlain()
	return s
}

func (s *Setting) AsDiscord() (*Discord, error) {
	return parseSettingAs(s, base.SettingTypeDiscord, func() *Discord { return &Discord{} })
}

func (s *Setting) MustAsDiscord() *Discord {
	return gofn.Must(s.AsDiscord())
}
