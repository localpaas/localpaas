package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
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
	Slack    *IMSlack    `json:"slack,omitempty"`
	Discord  *IMDiscord  `json:"discord,omitempty"`
	Telegram *IMTelegram `json:"telegram,omitempty"`
}

type IMSlack struct {
	Webhook EncryptedField `json:"webhook"`
}

type IMDiscord struct {
	Webhook EncryptedField `json:"webhook"`
}

type IMTelegram struct {
	BotToken EncryptedField `json:"botToken"`
	ChatID   string         `json:"chatId"`
}

func (s *IMService) GetType() base.SettingType {
	return base.SettingTypeIMService
}

func (s *IMService) GetRefObjectIDs() *RefObjectIDs {
	return &RefObjectIDs{}
}

func (s *IMService) CalcResLinks(setting *Setting) []*ResLink {
	return s.GetRefObjectIDs().CalcResLinks(base.ResourceTypeSetting, setting.ID)
}

func (s *IMService) Migrate(setting *Setting) (hasChange bool, err error) {
	if setting.Version == CurrentIMServiceVersion {
		return false, nil
	}
	if setting.Version > CurrentIMServiceVersion {
		return false, apperrors.New(apperrors.ErrDataVerNewerThanSystemVer)
	}

	// TODO: add migration if we make any change

	setting.Version = CurrentIMServiceVersion
	setting.UpdateVer++
	setting.MustSetData(s)
	return true, nil
}

func (s *IMService) MustDecrypt() *IMService {
	if s.Slack != nil {
		s.Slack.Webhook.MustGetPlain()
	}
	if s.Discord != nil {
		s.Discord.Webhook.MustGetPlain()
	}
	if s.Telegram != nil {
		s.Telegram.BotToken.MustGetPlain()
	}
	return s
}

func (s *Setting) AsIMService() (*IMService, error) {
	return parseSettingAs[*IMService](s)
}

func (s *Setting) MustAsIMService() *IMService {
	return gofn.Must(s.AsIMService())
}
