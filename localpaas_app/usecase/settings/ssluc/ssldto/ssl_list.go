package ssldto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ListSslReq struct {
	settings.ListSettingReq
}

func NewListSslReq() *ListSslReq {
	return &ListSslReq{
		ListSettingReq: settings.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListSslReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListSslResp struct {
	Meta *basedto.ListMeta `json:"meta"`
	Data []*SslResp        `json:"data"`
}

func TransformSsls(settings []*entity.Setting) (resp []*SslResp, err error) {
	resp = make([]*SslResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformSsl(setting)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
