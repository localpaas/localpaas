package githubappuc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/githubappuc/githubappdto"
)

const (
	currentSettingType    = base.SettingTypeGithubApp
	currentSettingVersion = entity.CurrentGithubAppVersion
)

func (uc *GithubAppUC) CreateGithubApp(
	ctx context.Context,
	auth *basedto.Auth,
	req *githubappdto.CreateGithubAppReq,
) (*githubappdto.CreateGithubAppResp, error) {
	req.Type = currentSettingType
	resp, err := providers.CreateSetting(ctx, uc.db, &req.CreateSettingReq, &providers.CreateSettingData{
		SettingRepo:   uc.settingRepo,
		VerifyingName: gofn.Coalesce(req.Name, req.Organization),
		Version:       currentSettingVersion,
		PrepareCreation: func(
			ctx context.Context,
			db database.Tx,
			data *providers.CreateSettingData,
			pData *providers.PersistingSettingCreationData,
		) error {
			pData.Setting.Kind = string(base.SettingTypeGithubApp)
			err := pData.Setting.SetData(&entity.GithubApp{
				ClientID:       req.ClientID,
				ClientSecret:   entity.NewEncryptedField(req.ClientSecret),
				Organization:   req.Organization,
				WebhookURL:     req.WebhookURL,
				WebhookSecret:  entity.NewEncryptedField(req.WebhookSecret),
				AppID:          req.GhAppID,
				InstallationID: req.GhInstallationID,
				PrivateKey:     entity.NewEncryptedField(req.PrivateKey),
				SSOEnabled:     req.SSOEnabled,
			})
			if err != nil {
				return apperrors.Wrap(err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &githubappdto.CreateGithubAppResp{
		Data: &githubappdto.GithubAppCreationResp{
			ID:          resp.Data.ID,
			CallbackURL: config.Current.SsoBaseCallbackURL() + "/" + resp.Data.ID,
		},
	}, nil
}
