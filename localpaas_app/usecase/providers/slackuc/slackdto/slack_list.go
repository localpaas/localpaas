package slackdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
)

type ListSlackReq struct {
	providers.ListSettingReq
}

func NewListSlackReq() *ListSlackReq {
	return &ListSlackReq{
		ListSettingReq: providers.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListSlackReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListSlackResp struct {
	Meta *basedto.Meta `json:"meta"`
	Data []*SlackResp  `json:"data"`
}

func TransformSlacks(settings []*entity.Setting) (resp []*SlackResp, err error) {
	resp = make([]*SlackResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformSlack(setting)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
