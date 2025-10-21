package projectdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type DeleteProjectEnvReq struct {
	ProjectID    string `json:"-"`
	ProjectEnvID string `json:"-"`
}

func NewDeleteProjectEnvReq() *DeleteProjectEnvReq {
	return &DeleteProjectEnvReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteProjectEnvReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.ProjectEnvID, true, "projectEnvId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteProjectEnvResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
