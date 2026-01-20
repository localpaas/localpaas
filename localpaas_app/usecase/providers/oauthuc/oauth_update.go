package oauthuc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/oauthuc/oauthdto"
)

func (uc *OAuthUC) UpdateOAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *oauthdto.UpdateOAuthReq,
) (*oauthdto.UpdateOAuthResp, error) {
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
			pData.Setting.Kind = gofn.Coalesce(string(req.Kind), pData.Setting.Kind)
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

	return &oauthdto.UpdateOAuthResp{}, nil
}
