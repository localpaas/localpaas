package webhookdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ListWebhookReq struct {
	settings.ListSettingReq
}

func NewListWebhookReq() *ListWebhookReq {
	return &ListWebhookReq{
		ListSettingReq: settings.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListWebhookReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListWebhookResp struct {
	Meta *basedto.ListMeta `json:"meta"`
	Data []*WebhookResp    `json:"data"`
}

func TransformWebhooks(settings []*entity.Setting) (resp []*WebhookResp, err error) {
	resp = make([]*WebhookResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformWebhook(setting)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
