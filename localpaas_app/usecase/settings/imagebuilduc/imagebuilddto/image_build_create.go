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
	NoCache   bool                    `json:"noCache"`
	NoVerbose bool                    `json:"noVerbose"`
}

func (req *ImageBuildBaseReq) ToEntity() *entity.ImageBuild {
	return &entity.ImageBuild{
		Resources: req.Resources.ToEntity(),
		NoCache:   req.NoCache,
		NoVerbose: req.NoVerbose,
	}
}

func (req *ImageBuildBaseReq) validate(field string) []vld.Validator {
	return nil
}

type ImageBuildResourcesReq struct {
	CPUs      int32 `json:"cpus"`
	MemMB     int64 `json:"memMB"`
	MemSwapMB int64 `json:"memSwapMB"`
	ShmSizeMB int64 `json:"shmSizeMB"`
}

func (req *ImageBuildResourcesReq) ToEntity() *entity.ImageBuildResources {
	return &entity.ImageBuildResources{
		CPUs:      req.CPUs,
		MemMB:     req.MemMB,
		MemSwapMB: req.MemSwapMB,
		ShmSizeMB: req.ShmSizeMB,
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
