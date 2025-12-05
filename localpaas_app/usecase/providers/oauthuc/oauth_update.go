package oauthuc

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
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/oauthuc/oauthdto"
)

func (uc *OAuthUC) UpdateOAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *oauthdto.UpdateOAuthReq,
) (*oauthdto.UpdateOAuthResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		oauthData := &updateOAuthData{}
		err := uc.loadOAuthDataForUpdate(ctx, db, req, oauthData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingOAuthData{}
		uc.prepareUpdatingOAuth(req.OAuthBaseReq, oauthData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &oauthdto.UpdateOAuthResp{}, nil
}

type updateOAuthData struct {
	Setting *entity.Setting
}

func (uc *OAuthUC) loadOAuthDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *oauthdto.UpdateOAuthReq,
	data *updateOAuthData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeOAuth, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if req.UpdateVer != setting.UpdateVer {
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}
	data.Setting = setting
	uc.preprocessRequest(base.OAuthType(setting.Kind), req.OAuthBaseReq)

	// If name changes, validate the new one
	if req.Organization != "" && !strings.EqualFold(setting.Name, req.Organization) {
		conflictSetting, _ := uc.settingRepo.GetByName(ctx, db, base.SettingTypeOAuth, req.Organization, false)
		if conflictSetting != nil {
			return apperrors.NewAlreadyExist("OAuth").
				WithMsgLog("oauth '%s' already exists", conflictSetting.Name)
		}
	}

	return nil
}

func (uc *OAuthUC) prepareUpdatingOAuth(
	req *oauthdto.OAuthBaseReq,
	data *updateOAuthData,
	persistingData *persistingOAuthData,
) {
	timeNow := timeutil.NowUTC()
	setting := data.Setting
	setting.Name = gofn.Coalesce(req.Organization, setting.Name)

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

	setting.UpdatedAt = timeNow
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
