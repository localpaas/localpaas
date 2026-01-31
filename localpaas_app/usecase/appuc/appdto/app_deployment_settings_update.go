package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

const (
	imageNameMaxLen = 100
)

type UpdateAppDeploymentSettingsReq struct {
	ProjectID string `json:"-"`
	AppID     string `json:"-"`

	ImageSource   *DeploymentImageSourceReq   `json:"imageSource"`
	RepoSource    *DeploymentRepoSourceReq    `json:"repoSource"`
	TarballSource *DeploymentTarballSourceReq `json:"tarballSource"`
	ActiveMethod  base.DeploymentMethod       `json:"activeMethod"`

	Command               *string `json:"command"`
	WorkingDir            *string `json:"workingDir"`
	PreDeploymentCommand  *string `json:"preDeploymentCommand"`
	PostDeploymentCommand *string `json:"postDeploymentCommand"`

	UpdateVer int `json:"updateVer"`
}

type DeploymentImageSourceReq struct {
	Enabled      bool                `json:"enabled"`
	Image        string              `json:"image"`
	RegistryAuth basedto.ObjectIDReq `json:"registryAuth"`
}

func (req *DeploymentImageSourceReq) validate(field string) (res []vld.Validator) {
	if req == nil {
		return
	}
	if field != "" {
		field += "."
	}
	res = append(res, basedto.ValidateStr(&req.Image, true, 1, imageNameMaxLen, field+"image")...)
	res = append(res, basedto.ValidateObjectIDReq(&req.RegistryAuth, false, field+"registryAuth")...)
	return res
}

type DeploymentRepoSourceReq struct {
	Enabled        bool                `json:"enabled"`
	BuildTool      base.BuildTool      `json:"buildTool"`
	RepoURL        string              `json:"repoUrl"`
	RepoRef        string              `json:"repoRef"` // can be branch name, tag...
	Credentials    basedto.ObjectIDReq `json:"credentials"`
	DockerfilePath string              `json:"dockerfilePath"` // for BuildToolDockerfile only
	ImageTags      []string            `json:"imageTags"`
	RegistryAuth   basedto.ObjectIDReq `json:"registryAuth"`
}

// nolint
func (req *DeploymentRepoSourceReq) validate(field string) (res []vld.Validator) {
	if req == nil {
		return
	}
	if field != "" {
		field += "."
	}
	return res
}

type DeploymentTarballSourceReq struct {
	Enabled bool `json:"enabled"`
	// TODO: add implementation
}

// nolint
func (req *DeploymentTarballSourceReq) validate(field string) (res []vld.Validator) {
	if req == nil {
		return
	}
	if field != "" {
		field += "."
	}
	return res
}

func NewUpdateAppDeploymentSettingsReq() *UpdateAppDeploymentSettingsReq {
	return &UpdateAppDeploymentSettingsReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateAppDeploymentSettingsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	validators = append(validators, req.ImageSource.validate("imageSource")...)
	validators = append(validators, req.RepoSource.validate("repoSource")...)
	validators = append(validators, req.TarballSource.validate("tarballSource")...)
	// TODO: add validation for deployment settings input
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateAppDeploymentSettingsResp struct {
	Meta *basedto.Meta                        `json:"meta"`
	Data *UpdateAppDeploymentSettingsDataResp `json:"data"`
}

type UpdateAppDeploymentSettingsDataResp struct {
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}
