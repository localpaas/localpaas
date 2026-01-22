package discorddto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ListDiscordReq struct {
	settings.ListSettingReq
}

func NewListDiscordReq() *ListDiscordReq {
	return &ListDiscordReq{
		ListSettingReq: settings.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListDiscordReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListDiscordResp struct {
	Meta *basedto.Meta  `json:"meta"`
	Data []*DiscordResp `json:"data"`
}

func TransformDiscords(settings []*entity.Setting, objectID string) (resp []*DiscordResp, err error) {
	resp = make([]*DiscordResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformDiscord(setting, objectID)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
