package oauthuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/oauthuc/oauthdto"
)

func (uc *OAuthUC) ListOAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *oauthdto.ListOAuthReq,
) (*oauthdto.ListOAuthResp, error) {
	req.Type = currentSettingType
	resp, err := settings.ListSetting(ctx, uc.db, auth, &req.ListSettingReq, &settings.ListSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := oauthdto.TransformOAuths(resp.Data, config.Current.SsoBaseCallbackURL())
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &oauthdto.ListOAuthResp{
		Meta: resp.Meta,
		Data: respData,
	}, nil
}
