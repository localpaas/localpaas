package registryauthdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type ListRegistryAuthReq struct {
	Status []base.SettingStatus `json:"-" mapstructure:"status"`
	Search string               `json:"-" mapstructure:"search"`

	Paging basedto.Paging `json:"-"`
}

func NewListRegistryAuthReq() *ListRegistryAuthReq {
	return &ListRegistryAuthReq{
		Paging: basedto.Paging{
			// Default paging if unset by client
			Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
		},
	}
}

func (req *ListRegistryAuthReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateSlice(req.Status, true, 0, base.AllSettingStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListRegistryAuthResp struct {
	Meta *basedto.Meta       `json:"meta"`
	Data []*RegistryAuthResp `json:"data"`
}

func TransformRegistryAuths(settings []*entity.Setting) (resp []*RegistryAuthResp, err error) {
	resp = make([]*RegistryAuthResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformRegistryAuth(setting)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
