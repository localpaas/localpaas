package secretuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/secretuc/secretdto"
)

func (uc *SecretUC) ListSecret(
	ctx context.Context,
	auth *basedto.Auth,
	req *secretdto.ListSecretReq,
) (*secretdto.ListSecretResp, error) {
	listOpts := []bunex.SelectQueryOption{
		bunex.SelectWhere("setting.type = ?", base.SettingTypeSecret),
	}

	if req.ObjectID != "" {
		listOpts = append(listOpts,
			bunex.SelectWhere("setting.object_id = ?", req.ObjectID))
	} else {
		listOpts = append(listOpts,
			bunex.SelectWhere("setting.object_id IS NULL"))
	}
	if len(req.Status) > 0 {
		listOpts = append(listOpts,
			bunex.SelectWhere("setting.status IN (?)", bunex.In(req.Status)))
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

	resp, err := secretdto.TransformSecrets(settings, false)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &secretdto.ListSecretResp{
		Meta: &basedto.Meta{Page: paging},
		Data: resp,
	}, nil
}
