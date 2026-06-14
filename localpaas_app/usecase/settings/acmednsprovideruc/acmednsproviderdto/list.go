package acmednsproviderdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ListAcmeDnsProviderReq struct {
	settings.ListSettingReq
}

func NewListAcmeDnsProviderReq() *ListAcmeDnsProviderReq {
	return &ListAcmeDnsProviderReq{
		ListSettingReq: settings.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListAcmeDnsProviderReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListAcmeDnsProviderResp struct {
	Meta *basedto.ListMeta      `json:"meta"`
	Data []*AcmeDnsProviderResp `json:"data"`
}

func TransformAcmeDnsProviders(
	settings []*entity.Setting,
	refObjects *entity.RefObjects,
) (resp []*AcmeDnsProviderResp, err error) {
	resp = make([]*AcmeDnsProviderResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformAcmeDnsProvider(setting, refObjects)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
