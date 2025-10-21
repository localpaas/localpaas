package userdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type ListUserBaseReq struct {
	Status []base.UserStatus `json:"-" mapstructure:"status"`
	Search string            `json:"-" mapstructure:"search"`

	Paging basedto.Paging `json:"-"`
}

func NewListUserBaseReq() *ListUserBaseReq {
	return &ListUserBaseReq{
		Status: []base.UserStatus{base.UserStatusActive},
		Paging: basedto.Paging{
			// Default paging if unset by client
			Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "full_name"}},
		},
	}
}

// Validate implements interface basedto.ReqValidator
func (req *ListUserBaseReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators,
		basedto.ValidateSlice(req.Status, true, 0, base.AllUserStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListUserBaseResp struct {
	Meta *basedto.Meta   `json:"meta"`
	Data []*UserBaseResp `json:"data"`
}
