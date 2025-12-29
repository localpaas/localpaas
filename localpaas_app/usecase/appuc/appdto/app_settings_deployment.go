package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
)

const (
	imageNameMaxLen = 100
)

//
// REQUEST
//

type DeploymentSettingsReq struct {
	ImageSource   *DeploymentImageSourceReq   `json:"imageSource"`
	RepoSource    *DeploymentRepoSourceReq    `json:"repoSource"`
	TarballSource *DeploymentTarballSourceReq `json:"tarballSource"`
	UpdateVer     int                         `json:"updateVer"`
}

func (req *DeploymentSettingsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ImageSource.validate("imageSource")...)
	validators = append(validators, req.RepoSource.validate("repoSource")...)
	validators = append(validators, req.TarballSource.validate("tarballSource")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

func (req *DeploymentSettingsReq) validate(_ string) []vld.Validator { //nolint
	if req == nil {
		return nil
	}
	// TODO: add validation
	return nil
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

//
// RESPONSE
//

type DeploymentSettingsResp struct {
	ImageSource   *DeploymentImageSourceResp   `json:"imageSource"`
	RepoSource    *DeploymentRepoSourceResp    `json:"repoSource"`
	TarballSource *DeploymentTarballSourceResp `json:"tarballSource"`
	UpdateVer     int                          `json:"updateVer"`
}

type DeploymentImageSourceResp struct {
	Enabled      bool                     `json:"enabled"`
	Image        string                   `json:"image"`
	RegistryAuth *basedto.NamedObjectResp `json:"registryAuth"`
}

type DeploymentRepoSourceResp struct {
	Enabled        bool                     `json:"enabled"`
	BuildTool      base.BuildTool           `json:"buildTool"`
	RepoURL        string                   `json:"repoUrl"`
	RepoRef        string                   `json:"repoRef"` // can be branch name, tag...
	Credentials    *RepoCredentialsResp     `json:"credentials"`
	DockerfilePath string                   `json:"dockerfilePath"` // for BuildToolDockerfile only
	ImageTags      []string                 `json:"imageTags"`
	RegistryAuth   *basedto.NamedObjectResp `json:"registryAuth"`
}

type RepoCredentialsResp struct {
	ID   string           `json:"id"`
	Type base.SettingType `json:"type"`
}

type DeploymentTarballSourceResp struct {
	Enabled bool `json:"enabled"`
}

func TransformDeploymentSettings(input *AppSettingsTransformationInput) (resp *DeploymentSettingsResp, err error) {
	if input.DeploymentSettings == nil {
		return nil, nil
	}
	if err = copier.Copy(&resp, input.DeploymentSettings); err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
