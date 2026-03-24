package cloudproviderdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ListCloudProviderReq struct {
	settings.ListSettingReq
}

func NewListCloudProviderReq() *ListCloudProviderReq {
	return &ListCloudProviderReq{
		ListSettingReq: settings.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListCloudProviderReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListCloudProviderResp struct {
	Meta *basedto.ListMeta    `json:"meta"`
	Data []*CloudProviderResp `json:"data"`
}

func TransformAWSs(
	settings []*entity.Setting,
	refObjects *entity.RefObjects,
) (resp []*CloudProviderResp, err error) {
	resp = make([]*CloudProviderResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformCloudProvider(setting, refObjects)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
