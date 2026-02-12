package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentIMServiceVersion = 1
)

var _ = registerSettingParser(base.SettingTypeIMService, &imServiceParser{})

type imServiceParser struct {
}

func (s *imServiceParser) New() SettingData {
	return &IMService{}
}

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

func (s *IMService) GetType() base.SettingType {
	return base.SettingTypeIMService
}

func (s *IMService) GetRefSettingIDs() []string {
	return nil
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
	return parseSettingAs[*IMService](s)
}

func (s *Setting) MustAsIMService() *IMService {
	return gofn.Must(s.AsIMService())
}
