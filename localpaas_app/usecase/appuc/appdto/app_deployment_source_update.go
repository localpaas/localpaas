package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

const (
	imageNameMaxLen = 100
)

type UpdateAppDeploymentSourceReq struct {
	ProjectID string `json:"-"`
	AppID     string `json:"-"`

	ImageSource *DeploymentImageSourceReq `json:"imageSource"`
	CodeSource  *DeploymentCodeSourceReq  `json:"codeSource"`
}

type DeploymentImageSourceReq struct {
	Enabled      bool                `json:"enabled"`
	Name         string              `json:"name"`
	RegistryAuth basedto.ObjectIDReq `json:"registryAuth"`
}

func (req *DeploymentImageSourceReq) validate(field string) (res []vld.Validator) {
	if req == nil {
		return
	}
	if field != "" {
		field += "."
	}
	res = append(res, basedto.ValidateStr(&req.Name, true, 1, imageNameMaxLen, field+"name")...)
	res = append(res, basedto.ValidateObjectIDReq(&req.RegistryAuth, false, field+"registryAuth")...)
	return res
}

type DeploymentCodeSourceReq struct {
	Enabled bool `json:"enabled"`
	// TODO: add implementation
}

// nolint
func (req *DeploymentCodeSourceReq) validate(field string) (res []vld.Validator) {
	if req == nil {
		return
	}
	if field != "" {
		field += "."
	}
	return res
}

func NewUpdateAppDeploymentSourceReq() *UpdateAppDeploymentSourceReq {
	return &UpdateAppDeploymentSourceReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateAppDeploymentSourceReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	validators = append(validators, req.ImageSource.validate("imageSource")...)
	validators = append(validators, req.CodeSource.validate("codeSource")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateAppDeploymentSourceResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
