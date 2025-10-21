package projectdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

const (
	minEnvNameLen = 1
	maxEnvNameLen = 100
)

type CreateProjectEnvReq struct {
	ProjectID string `json:"-"`
	Name      string `json:"name"`
}

func NewCreateProjectEnvReq() *CreateProjectEnvReq {
	return &CreateProjectEnvReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateProjectEnvReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateStr(&req.Name, true, minEnvNameLen, maxEnvNameLen, "tag")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateProjectEnvResp struct {
	Meta *basedto.BaseMeta     `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}
