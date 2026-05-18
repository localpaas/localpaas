package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentRepoWebhookVersion = 1
)

var _ = registerSettingParser(base.SettingTypeRepoWebhook, &repoWebhookParser{})

type repoWebhookParser struct {
}

func (s *repoWebhookParser) New() SettingData {
	return &RepoWebhook{}
}

type RepoWebhook struct {
	Kind   base.WebhookKind `json:"kind"`
	Secret string           `json:"secret"`
}

func (s *RepoWebhook) GetType() base.SettingType {
	return base.SettingTypeRepoWebhook
}

func (s *RepoWebhook) GetRefObjectIDs() *RefObjectIDs {
	return &RefObjectIDs{}
}

func (s *RepoWebhook) MustDecrypt() *RepoWebhook {
	return s
}

func (s *RepoWebhook) Migrate(setting *Setting) (hasChange bool, err error) {
	if setting.Version == CurrentRepoWebhookVersion {
		return false, nil
	}
	if setting.Version > CurrentRepoWebhookVersion {
		return false, apperrors.New(apperrors.ErrDataVerNewerThanSystemVer)
	}

	// TODO: add migration if we make any change

	setting.Version = CurrentRepoWebhookVersion
	setting.UpdateVer++
	setting.MustSetData(s)
	return true, nil
}

func (s *Setting) AsRepoWebhook() (*RepoWebhook, error) {
	// Github-app setting can be parsed as RepoWebhook
	if s.Type == base.SettingTypeGithubApp {
		ghApp, err := s.AsGithubApp()
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		return ghApp.ConvertAsRepoWebhook(), nil
	}
	return parseSettingAs[*RepoWebhook](s)
}

func (s *Setting) MustAsRepoWebhook() *RepoWebhook {
	return gofn.Must(s.AsRepoWebhook())
}
