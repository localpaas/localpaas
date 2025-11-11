package oauthuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/oauthuc/oauthdto"
)

func (uc *OAuthUC) GetOAuthNoAuth(
	ctx context.Context,
	req *oauthdto.GetOAuthNoAuthReq,
) (*oauthdto.OAuthResp, error) {
	listOpts := []bunex.SelectQueryOption{
		bunex.SelectWhere("setting.type = ?", base.SettingTypeOAuth),
		bunex.SelectLimit(1),
	}
	if req.ID != "" {
		listOpts = append(listOpts,
			bunex.SelectWhere("setting.id = ?", req.ID))
	}
	if req.Kind != "" {
		listOpts = append(listOpts,
			bunex.SelectWhere("setting.kind = ?", req.Kind))
	}
	if len(req.Status) > 0 {
		listOpts = append(listOpts,
			bunex.SelectWhere("setting.status IN (?)", bunex.In(req.Status)))
	}

	settings, _, err := uc.settingRepo.List(ctx, uc.db, nil, listOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if len(settings) == 0 {
		return nil, apperrors.NewNotFound("OAuth")
	}

	resp, err := oauthdto.TransformOAuth(settings[0], uc.GetOAuthBaseCallbackURL(), true)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return resp, nil
}
