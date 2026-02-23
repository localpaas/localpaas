package imagebuilddto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type CreateImageBuildReq struct {
	settings.CreateSettingReq
	*ImageBuildBaseReq
}

type ImageBuildBaseReq struct {
	Resources *ImageBuildResourcesReq `json:"resources"`
}

func (req *ImageBuildBaseReq) ToEntity() *entity.ImageBuild {
	return &entity.ImageBuild{
		Resources: req.Resources.ToEntity(),
	}
}

func (req *ImageBuildBaseReq) validate(field string) []vld.Validator {
	return nil
}

type ImageBuildResourcesReq struct {
	CPUs  uint32 `json:"cpus"`
	MemMB uint64 `json:"memMB"`
}

func (req *ImageBuildResourcesReq) ToEntity() *entity.ImageBuildResources {
	return &entity.ImageBuildResources{
		CPUs:  req.CPUs,
		MemMB: req.MemMB,
	}
}

func NewCreateImageBuildReq() *CreateImageBuildReq {
	return &CreateImageBuildReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateImageBuildReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateImageBuildResp struct {
	Meta *basedto.Meta         `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}
