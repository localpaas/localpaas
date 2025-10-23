package s3storagedto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type ListS3StorageReq struct {
	Search string `json:"-" mapstructure:"search"`

	Paging basedto.Paging `json:"-"`
}

func NewListS3StorageReq() *ListS3StorageReq {
	return &ListS3StorageReq{
		Paging: basedto.Paging{
			// Default paging if unset by client
			Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
		},
	}
}

func (req *ListS3StorageReq) Validate() apperrors.ValidationErrors {
	return nil
}

type ListS3StorageResp struct {
	Meta *basedto.Meta    `json:"meta"`
	Data []*S3StorageResp `json:"data"`
}

func TransformS3Storages(settings []*entity.Setting) ([]*S3StorageResp, error) {
	return basedto.TransformObjectSlice(settings, TransformS3Storage) //nolint:wrapcheck
}
