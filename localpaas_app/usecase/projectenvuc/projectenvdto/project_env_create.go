package projectenvdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type CreateProjectEnvReq struct {
	ProjectID string `json:"-"`
	*ProjectEnvReq
}

type ProjectEnvReq struct {
	Name   string             `json:"name"`
	Status base.ProjectStatus `json:"status"`
}

func (req *ProjectEnvReq) validate(field string) (res []vld.Validator) {
	if field != "" {
		field += "."
	}
	res = append(res, validateProjectEnvName(&req.Name, field+"name")...)
	res = append(res, basedto.ValidateStrIn(&req.Status, true, base.AllProjectStatuses, "status")...)
	return res
}

func NewCreateProjectEnvReq() *CreateProjectEnvReq {
	return &CreateProjectEnvReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateProjectEnvReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, req.ProjectEnvReq.validate("")...) //nolint
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateProjectEnvResp struct {
	Meta *basedto.BaseMeta     `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}
