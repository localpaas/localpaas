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
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/oauthuc/oauthdto"
)

var (
	// NOTE: only store special values
	mapNameByKind = map[string]string{
		string(base.OAuthTypeGitlabCustom): "Our Gitlab",
	}
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
	uc.preparePersistingOAuth(req.OAuthBaseReq, oauthData, persistingData)

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
	SettingKind string
}

func (uc *OAuthUC) loadOAuthData(
	ctx context.Context,
	db database.IDB,
	req *oauthdto.CreateOAuthReq,
	data *createOAuthData,
) error {
	uc.preprocessRequest(req.OAuthType, req.OAuthBaseReq)
	data.SettingKind = string(req.OAuthType)

	setting, err := uc.settingRepo.GetByName(ctx, db, base.SettingTypeOAuth, req.Name)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if setting != nil {
		return apperrors.NewAlreadyExist("OAuth").
			WithMsgLog("oauth setting '%s' already exists", req.Name)
	}

	return nil
}

func (uc *OAuthUC) preprocessRequest(
	oauthType base.OAuthType,
	req *oauthdto.OAuthBaseReq,
) {
	if !base.IsCustomOAuthType(oauthType) {
		req.Name = ""
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
) {
	timeNow := timeutil.NowUTC()
	setting := &entity.Setting{
		ID:        gofn.Must(ulid.NewStringULID()),
		Type:      base.SettingTypeOAuth,
		Status:    base.SettingStatusActive,
		Kind:      data.SettingKind,
		Name:      req.Name,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	if setting.Name == "" {
		setting.Name = mapNameByKind[data.SettingKind]
		if setting.Name == "" {
			setting.Name = gofn.StringToUpper1stLetter(setting.Kind)
		}
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
	setting.MustSetData(oauth.MustEncrypt())

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
