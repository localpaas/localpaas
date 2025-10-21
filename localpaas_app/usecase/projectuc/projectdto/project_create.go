package projectdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type CreateProjectReq struct {
	*ProjectBaseReq
}

type ProjectBaseReq struct {
	Name   string             `json:"name"`
	Status base.ProjectStatus `json:"status"`
	Tags   []string           `json:"tags"`
}

func NewCreateProjectReq() *CreateProjectReq {
	return &CreateProjectReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateProjectReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, validateProjectName(&req.Name, "name")...)
	validators = append(validators, basedto.ValidateStrIn(&req.Status, true, base.AllProjectStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateProjectResp struct {
	Meta *basedto.BaseMeta     `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}
