package githubappdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type ListGithubAppReq struct {
	Status []base.SettingStatus `json:"-" mapstructure:"status"`
	Search string               `json:"-" mapstructure:"search"`

	Paging basedto.Paging `json:"-"`
}

func NewListGithubAppReq() *ListGithubAppReq {
	return &ListGithubAppReq{
		Paging: basedto.Paging{
			// Default paging if unset by client
			Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
		},
	}
}

func (req *ListGithubAppReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateSlice(req.Status, true, 0, base.AllSettingStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListGithubAppResp struct {
	Meta *basedto.Meta    `json:"meta"`
	Data []*GithubAppResp `json:"data"`
}

func TransformGithubApps(settings []*entity.Setting, baseCallbackURL string) (resp []*GithubAppResp, err error) {
	resp = make([]*GithubAppResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformGithubApp(setting, baseCallbackURL)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
