package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
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
	ImageSource *DeploymentImageSourceReq `json:"imageSource"`
	CodeSource  *DeploymentCodeSourceReq  `json:"codeSource"`
}

func (req *DeploymentSettingsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ImageSource.validate("imageSource")...)
	validators = append(validators, req.CodeSource.validate("codeSource")...)
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

//
// RESPONSE
//

type DeploymentSettingsResp struct {
	ImageSource *DeploymentImageSourceResp `json:"imageSource"`
	CodeSource  *DeploymentCodeSourceResp  `json:"codeSource"`
}

type DeploymentImageSourceResp struct {
	Enabled      bool                `json:"enabled"`
	Name         string              `json:"name"`
	RegistryAuth basedto.ObjectIDReq `json:"registryAuth"`
}

type DeploymentCodeSourceResp struct {
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
