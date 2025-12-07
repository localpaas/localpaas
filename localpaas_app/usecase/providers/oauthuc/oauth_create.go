package oauthuc

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
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/oauthuc/oauthdto"
)

func (uc *OAuthUC) CreateOAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *oauthdto.CreateOAuthReq,
) (*oauthdto.CreateOAuthResp, error) {
	oauthData := &createOAuthData{}
	err := uc.loadOAuthData(ctx, uc.db, req, oauthData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	persistingData := &persistingOAuthData{}
	uc.preparePersistingOAuth(req, oauthData, persistingData)

	err = transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	createdItem := persistingData.UpsertingSettings[0]
	return &oauthdto.CreateOAuthResp{
		Data: &oauthdto.OAuthCreationResp{
			ID:          createdItem.ID,
			CallbackURL: config.Current.SsoBaseCallbackURL() + "/" + createdItem.ID,
		},
	}, nil
}

type createOAuthData struct {
}

func (uc *OAuthUC) loadOAuthData(
	ctx context.Context,
	db database.IDB,
	req *oauthdto.CreateOAuthReq,
	_ *createOAuthData,
) error {
	name := gofn.Coalesce(req.Name, req.Organization)
	setting, err := uc.settingRepo.GetByName(ctx, db, base.SettingTypeOAuth, name, false)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if setting != nil {
		return apperrors.NewAlreadyExist("OAuth").
			WithMsgLog("oauth setting '%s' already exists", name)
	}

	return nil
}

type persistingOAuthData struct {
	settingservice.PersistingSettingData
}

func (uc *OAuthUC) preparePersistingOAuth(
	req *oauthdto.CreateOAuthReq,
	_ *createOAuthData,
	persistingData *persistingOAuthData,
) {
	timeNow := timeutil.NowUTC()
	setting := &entity.Setting{
		ID:        gofn.Must(ulid.NewStringULID()),
		Type:      base.SettingTypeOAuth,
		Status:    base.SettingStatusActive,
		Kind:      string(req.Kind),
		Name:      gofn.Coalesce(req.Name, req.Organization),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	oauth := &entity.OAuth{
		ClientID:     req.ClientID,
		ClientSecret: entity.NewEncryptedField(req.ClientSecret),
		Organization: req.Organization,
		AuthURL:      req.AuthURL,
		TokenURL:     req.TokenURL,
		ProfileURL:   req.ProfileURL,
		Scopes:       req.Scopes,
	}
	setting.MustSetData(oauth)

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}

func (uc *OAuthUC) persistData(
	ctx context.Context,
	db database.IDB,
	persistingData *persistingOAuthData,
) error {
	err := uc.settingService.PersistSettingData(ctx, db, &persistingData.PersistingSettingData)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
