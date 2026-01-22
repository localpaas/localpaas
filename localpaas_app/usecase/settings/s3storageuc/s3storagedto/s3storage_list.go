package s3storagedto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ListS3StorageReq struct {
	settings.ListSettingReq
}

func NewListS3StorageReq() *ListS3StorageReq {
	return &ListS3StorageReq{
		ListSettingReq: settings.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListS3StorageReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListS3StorageResp struct {
	Meta *basedto.Meta    `json:"meta"`
	Data []*S3StorageResp `json:"data"`
}

func TransformS3Storages(settings []*entity.Setting) (resp []*S3StorageResp, err error) {
	resp = make([]*S3StorageResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformS3Storage(setting)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
