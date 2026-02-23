package imagebuilddto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type GetImageBuildReq struct {
	settings.GetSettingReq
}

func NewGetImageBuildReq() *GetImageBuildReq {
	return &GetImageBuildReq{}
}

func (req *GetImageBuildReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetImageBuildResp struct {
	Meta *basedto.Meta   `json:"meta"`
	Data *ImageBuildResp `json:"data"`
}

type ImageBuildResp struct {
	*settings.BaseSettingResp
	Resources *ImageBuildResourcesResp `json:"resources"`
}

type ImageBuildResourcesResp struct {
	CPUs  int `json:"cpus"`
	MemMB int `json:"memMB"`
}

func TransformImageBuild(
	setting *entity.Setting,
	_ *entity.RefObjects,
) (resp *ImageBuildResp, err error) {
	config := setting.MustAsImageBuild()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return resp, nil
}
