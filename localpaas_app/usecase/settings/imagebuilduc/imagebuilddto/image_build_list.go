package imagebuilddto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ListImageBuildReq struct {
	settings.ListSettingReq
}

func NewListImageBuildReq() *ListImageBuildReq {
	return &ListImageBuildReq{
		ListSettingReq: settings.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListImageBuildReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListImageBuildResp struct {
	Meta *basedto.ListMeta `json:"meta"`
	Data []*ImageBuildResp `json:"data"`
}

func TransformImageBuilds(
	settings []*entity.Setting,
	refObjects *entity.RefObjects,
) (resp []*ImageBuildResp, err error) {
	resp = make([]*ImageBuildResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformImageBuild(setting, refObjects)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
