package s3storageuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/s3storageuc/s3storagedto"
)

func (uc *S3StorageUC) ListS3Storage(
	ctx context.Context,
	auth *basedto.Auth,
	req *s3storagedto.ListS3StorageReq,
) (*s3storagedto.ListS3StorageResp, error) {
	listOpts := []bunex.SelectQueryOption{}

	if req.Search != "" {
		keyword := bunex.MakeLikeOpStr(req.Search, true)
		listOpts = append(listOpts,
			bunex.SelectWhereGroup(
				bunex.SelectWhere("s3_storage.name ILIKE ?", keyword),
			),
		)
	}

	s3Storages, paging, err := uc.s3StorageRepo.List(ctx, uc.db, &req.Paging, listOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := s3storagedto.TransformS3Storages(s3Storages)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &s3storagedto.ListS3StorageResp{
		Meta: &basedto.Meta{Page: paging},
		Data: resp,
	}, nil
}
