package awsdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ListAWSReq struct {
	settings.ListSettingReq
}

func NewListAWSReq() *ListAWSReq {
	return &ListAWSReq{
		ListSettingReq: settings.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListAWSReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListAWSResp struct {
	Meta *basedto.ListMeta `json:"meta"`
	Data []*AWSResp        `json:"data"`
}

func TransformAWSs(settings []*entity.Setting) (resp []*AWSResp, err error) {
	resp = make([]*AWSResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformAWS(setting)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
