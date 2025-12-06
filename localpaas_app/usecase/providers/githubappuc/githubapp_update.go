package githubappuc

import (
	"context"
	"strings"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/githubappuc/githubappdto"
)

func (uc *GithubAppUC) UpdateGithubApp(
	ctx context.Context,
	auth *basedto.Auth,
	req *githubappdto.UpdateGithubAppReq,
) (*githubappdto.UpdateGithubAppResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		appData := &updateGithubAppData{}
		err := uc.loadGithubAppDataForUpdate(ctx, db, req, appData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingGithubAppData{}
		uc.prepareUpdatingGithubApp(req.GithubAppBaseReq, appData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &githubappdto.UpdateGithubAppResp{}, nil
}

type updateGithubAppData struct {
	Setting *entity.Setting
}

func (uc *GithubAppUC) loadGithubAppDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *githubappdto.UpdateGithubAppReq,
	data *updateGithubAppData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeGithubApp, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if req.UpdateVer != setting.UpdateVer {
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}
	data.Setting = setting

	// If name changes, validate the new one
	if req.Organization != "" && !strings.EqualFold(setting.Name, req.Organization) {
		conflictSetting, _ := uc.settingRepo.GetByName(ctx, db, base.SettingTypeGithubApp, req.Organization, false)
		if conflictSetting != nil {
			return apperrors.NewAlreadyExist("GithubApp").
				WithMsgLog("github app '%s' already exists", conflictSetting.Name)
		}
	}

	return nil
}

func (uc *GithubAppUC) prepareUpdatingGithubApp(
	req *githubappdto.GithubAppBaseReq,
	data *updateGithubAppData,
	persistingData *persistingGithubAppData,
) {
	timeNow := timeutil.NowUTC()
	setting := data.Setting
	setting.Name = gofn.Coalesce(req.Name, req.Organization, setting.Name)

	githubApp := &entity.GithubApp{
		ClientID:       req.ClientID,
		ClientSecret:   entity.NewEncryptedField(req.ClientSecret),
		Organization:   req.Organization,
		WebhookURL:     req.WebhookURL,
		WebhookSecret:  entity.NewEncryptedField(req.WebhookSecret),
		AppID:          req.AppID,
		InstallationID: req.InstallationID,
		PrivateKey:     entity.NewEncryptedField(req.PrivateKey),
		SSOEnabled:     req.SSOEnabled,
	}
	setting.MustSetData(githubApp)

	setting.UpdatedAt = timeNow
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
