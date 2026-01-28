package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentIMServiceVersion = 1
)

type IMService struct {
	Slack   *Slack   `json:"slack,omitempty"`
	Discord *Discord `json:"discord,omitempty"`
}

type Slack struct {
	Webhook EncryptedField `json:"webhook"`
}

type Discord struct {
	Webhook EncryptedField `json:"webhook"`
}

func (s *IMService) MustDecrypt() *IMService {
	if s.Slack != nil {
		s.Slack.Webhook.MustGetPlain()
	}
	if s.Discord != nil {
		s.Discord.Webhook.MustGetPlain()
	}
	return s
}

func (s *Setting) AsIMService() (*IMService, error) {
	return parseSettingAs(s, base.SettingTypeIMService, func() *IMService { return &IMService{} })
}

func (s *Setting) MustAsIMService() *IMService {
	return gofn.Must(s.AsIMService())
}
