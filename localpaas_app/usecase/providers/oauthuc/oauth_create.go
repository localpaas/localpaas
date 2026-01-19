package oauthuc

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
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/oauthuc/oauthdto"
)

const (
	currentSettingType    = base.SettingTypeOAuth
	currentSettingVersion = entity.CurrentOAuthVersion
)

func (uc *OAuthUC) CreateOAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *oauthdto.CreateOAuthReq,
) (*oauthdto.CreateOAuthResp, error) {
	req.Type = currentSettingType
	resp, err := providers.CreateSetting(ctx, uc.db, &req.CreateSettingReq, &providers.CreateSettingData{
		SettingRepo:   uc.settingRepo,
		VerifyingName: gofn.Coalesce(req.Name, req.Organization),
		Version:       currentSettingVersion,
		PrepareCreation: func(ctx context.Context, db database.Tx, data *providers.CreateSettingData,
			pData *providers.PersistingSettingCreationData) error {
			pData.Setting.Kind = string(req.Kind)
			err := pData.Setting.SetData(&entity.OAuth{
				ClientID:     req.ClientID,
				ClientSecret: entity.NewEncryptedField(req.ClientSecret),
				Organization: req.Organization,
				AuthURL:      req.AuthURL,
				TokenURL:     req.TokenURL,
				ProfileURL:   req.ProfileURL,
				Scopes:       req.Scopes,
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

	return &oauthdto.CreateOAuthResp{
		Data: &oauthdto.OAuthCreationResp{
			ID:          resp.Data.ID,
			CallbackURL: config.Current.SsoBaseCallbackURL() + "/" + resp.Data.ID,
		},
	}, nil
}
