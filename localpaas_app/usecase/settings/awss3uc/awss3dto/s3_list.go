package awss3dto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ListAWSS3Req struct {
	settings.ListSettingReq
}

func NewListAWSS3Req() *ListAWSS3Req {
	return &ListAWSS3Req{
		ListSettingReq: settings.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListAWSS3Req) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListAWSS3Resp struct {
	Meta *basedto.ListMeta `json:"meta"`
	Data []*AWSS3Resp      `json:"data"`
}

func TransformAWSS3s(settings []*entity.Setting) (resp []*AWSS3Resp, err error) {
	resp = make([]*AWSS3Resp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformAWSS3(setting)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
