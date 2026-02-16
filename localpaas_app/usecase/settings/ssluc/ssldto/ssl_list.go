package ssldto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ListSSLReq struct {
	settings.ListSettingReq
}

func NewListSSLReq() *ListSSLReq {
	return &ListSSLReq{
		ListSettingReq: settings.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListSSLReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListSSLResp struct {
	Meta *basedto.ListMeta `json:"meta"`
	Data []*SSLResp        `json:"data"`
}

func TransformSSLs(
	settings []*entity.Setting,
	refObjects *entity.RefObjects,
) (resp []*SSLResp, err error) {
	resp = make([]*SSLResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformSSL(setting, refObjects)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
