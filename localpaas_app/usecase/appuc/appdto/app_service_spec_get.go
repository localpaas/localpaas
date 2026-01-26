package appdto

import (
	"github.com/docker/docker/api/types/swarm"
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/services/docker"
)

type GetAppServiceSpecReq struct {
	ProjectID string `json:"-"`
	AppID     string `json:"-"`
}

func NewGetAppServiceSpecReq() *GetAppServiceSpecReq {
	return &GetAppServiceSpecReq{}
}

func (req *GetAppServiceSpecReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetAppServiceSpecResp struct {
	Meta *basedto.Meta    `json:"meta"`
	Data *ServiceSpecResp `json:"data"`
}

type ServiceSpecResp struct {
	ServiceMode   *docker.ServiceModeSpec `json:"serviceMode"`
	EndpointSpec  *docker.EndpointSpec    `json:"endpointSpec"`
	TaskSpec      *docker.TaskSpec        `json:"taskSpec"`
	ContainerSpec *docker.ContainerSpec   `json:"containerSpec"`

	UpdateVer int `json:"updateVer"`
}

func TransformAppServiceSpec(service *swarm.Service) (resp *ServiceSpecResp, err error) {
	spec := &service.Spec
	resp = &ServiceSpecResp{
		ServiceMode:   docker.ConvertFromServiceModeSpec(spec),
		EndpointSpec:  docker.ConvertFromServiceEndpointSpec(spec.EndpointSpec),
		TaskSpec:      docker.ConvertFromServiceTaskSpec(&spec.TaskTemplate),
		ContainerSpec: docker.ConvertFromServiceContainerSpec(spec.TaskTemplate.ContainerSpec),
		UpdateVer:     int(service.Version.Index), //nolint:gosec
	}
	return resp, nil
}
