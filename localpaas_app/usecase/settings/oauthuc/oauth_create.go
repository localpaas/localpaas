package oauthuc

import (
	"context"
	"errors"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/oauthuc/oauthdto"
	"github.com/localpaas/localpaas/pkg/timeutil"
	"github.com/localpaas/localpaas/pkg/ulid"
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
	err = uc.preparePersistingOAuth(req.OAuthBaseReq, oauthData, persistingData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	err = transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	createdItem := persistingData.UpsertingSettings[0]
	return &oauthdto.CreateOAuthResp{
		Data: &basedto.ObjectIDResp{ID: createdItem.ID},
	}, nil
}

type createOAuthData struct {
	SettingName string
}

func (uc *OAuthUC) loadOAuthData(
	ctx context.Context,
	db database.IDB,
	req *oauthdto.CreateOAuthReq,
	data *createOAuthData,
) error {
	uc.preprocessRequest(req.OAuthType, req.OAuthBaseReq)

	settingName := string(req.OAuthType)
	data.SettingName = settingName

	setting, err := uc.settingRepo.GetByName(ctx, db, base.SettingTypeOAuth, settingName)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if setting != nil {
		return apperrors.NewAlreadyExist("OAuth").
			WithMsgLog("oauth setting '%s' already exists", settingName)
	}

	return nil
}

func (uc *OAuthUC) preprocessRequest(
	oauthType base.OAuthType,
	req *oauthdto.OAuthBaseReq,
) {
	if !base.IsCustomOAuthType(oauthType) {
		req.CallbackURL = ""
		req.AuthURL = ""
		req.TokenURL = ""
		req.ProfileURL = ""
	}
}

type persistingOAuthData struct {
	settingservice.PersistingSettingData
}

func (uc *OAuthUC) preparePersistingOAuth(
	req *oauthdto.OAuthBaseReq,
	data *createOAuthData,
	persistingData *persistingOAuthData,
) (err error) {
	timeNow := timeutil.NowUTC()
	setting := &entity.Setting{
		ID:        gofn.Must(ulid.NewStringULID()),
		Type:      base.SettingTypeOAuth,
		Status:    base.SettingStatusActive,
		Name:      data.SettingName,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	oauth := &entity.OAuth{
		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,
		Organization: req.Organization,
		CallbackURL:  req.CallbackURL,
		AuthURL:      req.AuthURL,
		TokenURL:     req.TokenURL,
		ProfileURL:   req.ProfileURL,
		Scopes:       req.Scopes,
	}
	err = setting.SetData(oauth)
	if err != nil {
		return apperrors.Wrap(err)
	}

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)

	return nil
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
