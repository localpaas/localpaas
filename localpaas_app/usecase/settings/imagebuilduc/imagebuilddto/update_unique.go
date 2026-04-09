package imagebuilddto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateUniqueImageBuildReq struct {
	settings.UpdateUniqueSettingReq
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

// nolint
func (req *ImageBuildBaseReq) validate(field string) (res []vld.Validator) {
	if field != "" {
		field += "."
	}
	// TODO: add validation
	return res
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

func NewUpdateUniqueImageBuildReq() *UpdateUniqueImageBuildReq {
	return &UpdateUniqueImageBuildReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateUniqueImageBuildReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateUniqueImageBuildResp struct {
	Meta *basedto.Meta `json:"meta"`
}
