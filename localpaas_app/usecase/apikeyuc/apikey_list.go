package apikeyuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/apikeyuc/apikeydto"
)

func (uc *APIKeyUC) ListAPIKey(
	ctx context.Context,
	auth *basedto.Auth,
	req *apikeydto.ListAPIKeyReq,
) (*apikeydto.ListAPIKeyResp, error) {
	listOpts := []bunex.SelectQueryOption{
		bunex.SelectWhere("setting.deleted_at IS NULL"),
		bunex.SelectWhere("setting.type = ?", base.SettingTypeAPIKey),
		bunex.SelectWhere("setting.object_id = ?", auth.User.ID),
		bunex.SelectRelation("ObjectUser", bunex.SelectWithDeleted()),
	}
	if len(req.Status) > 0 {
		listOpts = append(listOpts, bunex.SelectWhere("setting.status IN (?)", bunex.In(req.Status)))
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

	resp, err := apikeydto.TransformAPIKeys(settings)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &apikeydto.ListAPIKeyResp{
		Meta: &basedto.Meta{Page: paging},
		Data: resp,
	}, nil
}
