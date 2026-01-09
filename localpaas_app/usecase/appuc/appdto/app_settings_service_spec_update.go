package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/services/docker"
)

type UpdateAppServiceSpecReq struct {
	ProjectID string `json:"-"`
	AppID     string `json:"-"`

	ServiceMode   *docker.ServiceModeSpec `json:"serviceMode"`
	EndpointSpec  *docker.EndpointSpec    `json:"endpointSpec"`
	TaskSpec      *docker.TaskSpec        `json:"taskSpec"`
	ContainerSpec *docker.ContainerSpec   `json:"containerSpec"`

	UpdateVer int `json:"updateVer"`
}

func NewUpdateAppServiceSpecReq() *UpdateAppServiceSpecReq {
	return &UpdateAppServiceSpecReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateAppServiceSpecReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	// TODO: validate env var input
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateAppServiceSpecResp struct {
	Meta *basedto.BaseMeta             `json:"meta"`
	Data *UpdateAppServiceSpecDataResp `json:"data"`
}

type UpdateAppServiceSpecDataResp struct {
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}
