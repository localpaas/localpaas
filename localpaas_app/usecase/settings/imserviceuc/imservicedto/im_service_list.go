package imservicedto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ListIMServiceReq struct {
	settings.ListSettingReq
}

func NewListIMServiceReq() *ListIMServiceReq {
	return &ListIMServiceReq{
		ListSettingReq: settings.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListIMServiceReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListIMServiceResp struct {
	Meta *basedto.ListMeta `json:"meta"`
	Data []*IMServiceResp  `json:"data"`
}

func TransformIMServices(
	settings []*entity.Setting,
	refObjects *entity.RefObjects,
) (resp []*IMServiceResp, err error) {
	resp = make([]*IMServiceResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformIMService(setting, refObjects)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
