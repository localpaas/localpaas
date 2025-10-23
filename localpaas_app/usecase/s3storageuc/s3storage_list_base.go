package s3storageuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/s3storageuc/s3storagedto"
)

func (uc *S3StorageUC) ListS3StorageBase(
	ctx context.Context,
	auth *basedto.Auth,
	req *s3storagedto.ListS3StorageBaseReq,
) (*s3storagedto.ListS3StorageBaseResp, error) {
	listOpts := []bunex.SelectQueryOption{
		bunex.SelectWhere("setting.type = ?", base.SettingTypeS3Storage),
	}

	if req.Search != "" {
		keyword := bunex.MakeLikeOpStr(req.Search, true)
		listOpts = append(listOpts,
			bunex.SelectWhereGroup(
				bunex.SelectWhere("setting.name ILIKE ?", keyword),
			),
		)
	}

	settings, pagingMeta, err := uc.settingRepo.List(ctx, uc.db, &req.Paging, listOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := s3storagedto.TransformS3StoragesBase(settings)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &s3storagedto.ListS3StorageBaseResp{
		Meta: &basedto.Meta{Page: pagingMeta},
		Data: resp,
	}, nil
}
