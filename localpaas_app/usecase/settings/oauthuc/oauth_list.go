package oauthuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/oauthuc/oauthdto"
)

func (uc *OAuthUC) ListOAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *oauthdto.ListOAuthReq,
) (*oauthdto.ListOAuthResp, error) {
	listOpts := []bunex.SelectQueryOption{
		bunex.SelectWhere("setting.type = ?", base.SettingTypeOAuth),
	}

	if req.Search != "" {
		keyword := bunex.MakeLikeOpStr(req.Search, true)
		listOpts = append(listOpts,
			bunex.SelectWhereGroup(
				bunex.SelectWhere("setting.name ILIKE ?", keyword),
			),
		)
	}

	settings, paging, err := uc.settingRepo.List(ctx, uc.db, &req.Paging, listOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := oauthdto.TransformOAuths(settings)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &oauthdto.ListOAuthResp{
		Meta: &basedto.Meta{Page: paging},
		Data: resp,
	}, nil
}
