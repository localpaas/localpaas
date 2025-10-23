package s3storagedto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type ListS3StorageBaseReq struct {
	Search string `json:"-" mapstructure:"search"`

	Paging basedto.Paging `json:"-"`
}

func NewListS3StorageBaseReq() *ListS3StorageBaseReq {
	return &ListS3StorageBaseReq{
		Paging: basedto.Paging{
			// Default paging if unset by client
			Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
		},
	}
}

// Validate implements interface basedto.ReqValidator
func (req *ListS3StorageBaseReq) Validate() apperrors.ValidationErrors {
	return nil
}

type ListS3StorageBaseResp struct {
	Meta *basedto.Meta        `json:"meta"`
	Data []*S3StorageBaseResp `json:"data"`
}
