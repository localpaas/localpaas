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

func TransformGitTokens(settings []*entity.Setting) ([]*GitTokenResp, error) {
	resp, err := basedto.TransformObjectSlice(settings, TransformGitToken)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
