package githubappuc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/githubappuc/githubappdto"
)

func (uc *GithubAppUC) UpdateGithubApp(
	ctx context.Context,
	auth *basedto.Auth,
	req *githubappdto.UpdateGithubAppReq,
) (*githubappdto.UpdateGithubAppResp, error) {
	req.Type = currentSettingType
	_, err := providers.UpdateSetting(ctx, uc.db, &req.UpdateSettingReq, &providers.UpdateSettingData{
		SettingRepo:   uc.settingRepo,
		VerifyingName: req.Name,
		PrepareUpdate: func(
			ctx context.Context,
			db database.Tx,
			data *providers.UpdateSettingData,
			pData *providers.PersistingSettingData,
		) error {
			pData.Setting.Name = gofn.Coalesce(req.Name, pData.Setting.Name)
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

	return &githubappdto.UpdateGithubAppResp{}, nil
}
