package gittokendto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ListGitTokenReq struct {
	settings.ListSettingReq
}

func NewListGitTokenReq() *ListGitTokenReq {
	return &ListGitTokenReq{
		ListSettingReq: settings.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListGitTokenReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListGitTokenResp struct {
	Meta *basedto.Meta   `json:"meta"`
	Data []*GitTokenResp `json:"data"`
}

func TransformGitTokens(settings []*entity.Setting, objectID string) (resp []*GitTokenResp, err error) {
	resp = make([]*GitTokenResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformGitToken(setting, objectID)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
