package oauthuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/oauthuc/oauthdto"
)

func (uc *OAuthUC) GetOAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *oauthdto.GetOAuthReq,
) (*oauthdto.GetOAuthResp, error) {
	setting, err := uc.settingRepo.GetByID(ctx, uc.db, req.ID,
		bunex.SelectWhere("setting.type = ?", base.SettingTypeOAuth),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := oauthdto.TransformOAuth(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &oauthdto.GetOAuthResp{
		Data: resp,
	}, nil
}
