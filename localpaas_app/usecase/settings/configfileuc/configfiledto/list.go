package configfiledto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ListConfigFileReq struct {
	settings.ListSettingReq
}

func NewListConfigFileReq() *ListConfigFileReq {
	return &ListConfigFileReq{
		ListSettingReq: settings.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListConfigFileReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListConfigFileResp struct {
	Meta *basedto.ListMeta `json:"meta"`
	Data []*ConfigFileResp `json:"data"`
}

func TransformConfigFiles(
	settings []*entity.Setting,
	refObjects *entity.RefObjects,
) (resp []*ConfigFileResp, err error) {
	resp = make([]*ConfigFileResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformConfigFile(setting, refObjects)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
