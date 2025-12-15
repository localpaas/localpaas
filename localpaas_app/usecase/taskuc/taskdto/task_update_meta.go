package taskdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateTaskMetaReq struct {
	ID        string           `json:"-"`
	Status    *base.TaskStatus `json:"status"`
	UpdateVer int              `json:"updateVer"`
}

func NewUpdateTaskMetaReq() *UpdateTaskMetaReq {
	return &UpdateTaskMetaReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateTaskMetaReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStrIn(req.Status, false,
		base.AllTaskSettableStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateTaskMetaResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
