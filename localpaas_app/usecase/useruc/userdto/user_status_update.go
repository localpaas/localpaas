package userdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateStatusReq struct {
	ID     string          `json:"-"`
	Status base.UserStatus `json:"status"`
}

func NewUpdateStatusReq() *UpdateStatusReq {
	return &UpdateStatusReq{}
}

func (req *UpdateStatusReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "status")...)
	validators = append(validators, basedto.ValidateStrIn(&req.Status, true, base.AllUserStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateStatusResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
