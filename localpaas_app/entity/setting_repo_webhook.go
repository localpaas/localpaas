package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentRepoWebhookVersion = 1
)

type RepoWebhook struct {
	Kind   base.WebhookKind `json:"kind"`
	Secret string           `json:"secret"`
}

func (s *RepoWebhook) GetType() base.SettingType {
	return base.SettingTypeRepoWebhook
}

func (s *RepoWebhook) GetRefSettingIDs() []string {
	return nil
}

func (s *RepoWebhook) MustDecrypt() *RepoWebhook {
	return s
}

func (s *Setting) AsRepoWebhook() (*RepoWebhook, error) {
	return parseSettingAs(s, func() *RepoWebhook { return &RepoWebhook{} })
}

func (s *Setting) MustAsRepoWebhook() *RepoWebhook {
	return gofn.Must(s.AsRepoWebhook())
}
