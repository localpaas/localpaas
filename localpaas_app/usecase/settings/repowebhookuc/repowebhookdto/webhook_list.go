package repowebhookdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ListRepoWebhookReq struct {
	settings.ListSettingReq
}

func NewListRepoWebhookReq() *ListRepoWebhookReq {
	return &ListRepoWebhookReq{
		ListSettingReq: settings.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListRepoWebhookReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListRepoWebhookResp struct {
	Meta *basedto.ListMeta  `json:"meta"`
	Data []*RepoWebhookResp `json:"data"`
}

func TransformRepoWebhooks(settings []*entity.Setting) (resp []*RepoWebhookResp, err error) {
	resp = make([]*RepoWebhookResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformRepoWebhook(setting)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
