package apikeydto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ListAPIKeyReq struct {
	settings.ListSettingReq
}

func NewListAPIKeyReq() *ListAPIKeyReq {
	return &ListAPIKeyReq{
		ListSettingReq: settings.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListAPIKeyReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListAPIKeyResp struct {
	Meta *basedto.ListMeta `json:"meta"`
	Data []*APIKeyResp     `json:"data"`
}

func TransformAPIKeys(settings []*entity.Setting) (resp []*APIKeyResp, err error) {
	resp = make([]*APIKeyResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformAPIKey(setting)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
