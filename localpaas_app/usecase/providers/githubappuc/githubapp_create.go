package githubappuc

import (
	"context"
	"errors"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/githubappuc/githubappdto"
)

func (uc *GithubAppUC) CreateGithubApp(
	ctx context.Context,
	auth *basedto.Auth,
	req *githubappdto.CreateGithubAppReq,
) (*githubappdto.CreateGithubAppResp, error) {
	appData := &createGithubAppData{}
	err := uc.loadGithubAppData(ctx, uc.db, req, appData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	persistingData := &persistingGithubAppData{}
	uc.preparePersistingGithubApp(req, appData, persistingData)

	err = transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	createdItem := persistingData.UpsertingSettings[0]
	return &githubappdto.CreateGithubAppResp{
		Data: &githubappdto.GithubAppCreationResp{
			ID:          createdItem.ID,
			CallbackURL: config.Current.SsoBaseCallbackURL() + "/" + createdItem.ID,
		},
	}, nil
}

type createGithubAppData struct {
}

func (uc *GithubAppUC) loadGithubAppData(
	ctx context.Context,
	db database.IDB,
	req *githubappdto.CreateGithubAppReq,
	_ *createGithubAppData,
) error {
	name := gofn.Coalesce(req.Name, req.Organization)
	setting, err := uc.settingRepo.GetByName(ctx, db, base.SettingTypeGithubApp, name, false)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if setting != nil {
		return apperrors.NewAlreadyExist("GithubApp").
			WithMsgLog("github app setting '%s' already exists", name)
	}

	return nil
}

type persistingGithubAppData struct {
	settingservice.PersistingSettingData
}

func (uc *GithubAppUC) preparePersistingGithubApp(
	req *githubappdto.CreateGithubAppReq,
	_ *createGithubAppData,
	persistingData *persistingGithubAppData,
) {
	timeNow := timeutil.NowUTC()
	setting := &entity.Setting{
		ID:        gofn.Must(ulid.NewStringULID()),
		Type:      base.SettingTypeGithubApp,
		Status:    base.SettingStatusActive,
		Kind:      string(base.SettingTypeGithubApp),
		Name:      gofn.Coalesce(req.Name, req.Organization),
		Version:   entity.CurrentGithubAppVersion,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

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

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}

func (uc *GithubAppUC) persistData(
	ctx context.Context,
	db database.IDB,
	persistingData *persistingGithubAppData,
) error {
	err := uc.settingService.PersistSettingData(ctx, db, &persistingData.PersistingSettingData)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
