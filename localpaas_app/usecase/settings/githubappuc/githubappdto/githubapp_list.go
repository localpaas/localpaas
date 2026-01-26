package githubappdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ListGithubAppReq struct {
	settings.ListSettingReq
}

func NewListGithubAppReq() *ListGithubAppReq {
	return &ListGithubAppReq{
		ListSettingReq: settings.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListGithubAppReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListGithubAppResp struct {
	Meta *basedto.ListMeta `json:"meta"`
	Data []*GithubAppResp  `json:"data"`
}

func TransformGithubApps(settings []*entity.Setting, baseCallbackURL string, objectID string) (
	resp []*GithubAppResp, err error) {
	resp = make([]*GithubAppResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformGithubApp(setting, baseCallbackURL, objectID)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
