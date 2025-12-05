package oauthuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/oauthuc/oauthdto"
)

func (uc *OAuthUC) GetOAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *oauthdto.GetOAuthReq,
) (*oauthdto.GetOAuthResp, error) {
	setting, err := uc.settingRepo.GetByID(ctx, uc.db, base.SettingTypeOAuth, req.ID, false)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	setting.MustAsOAuth().MustDecrypt()
	resp, err := oauthdto.TransformOAuth(setting, config.Current.SsoBaseCallbackURL())
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &oauthdto.GetOAuthResp{
		Data: resp,
	}, nil
}
