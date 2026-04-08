package cloudstoragedto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ListCloudStorageReq struct {
	settings.ListSettingReq
}

func NewListCloudStorageReq() *ListCloudStorageReq {
	return &ListCloudStorageReq{
		ListSettingReq: settings.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListCloudStorageReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListCloudStorageResp struct {
	Meta *basedto.ListMeta   `json:"meta"`
	Data []*CloudStorageResp `json:"data"`
}

func TransformCloudStorages(
	settings []*entity.Setting,
	refObjects *entity.RefObjects,
) (resp []*CloudStorageResp, err error) {
	resp = make([]*CloudStorageResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformCloudStorage(setting, refObjects)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
