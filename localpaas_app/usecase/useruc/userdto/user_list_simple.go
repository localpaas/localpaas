package userdto

import (
	vld "github.com/tiendc/go-validator"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type ListUserSimpleReq struct {
	Status []base.UserStatus `json:"-" mapstructure:"status"`
	Search string            `json:"-" mapstructure:"search"`

	Paging basedto.Paging `json:"-"`
}

func NewListUserSimpleReq() *ListUserSimpleReq {
	return &ListUserSimpleReq{
		Status: []base.UserStatus{base.UserStatusActive},
		Paging: basedto.Paging{
			// Default paging if unset by client
			Sort: basedto.Orders{
				{Direction: basedto.DirectionAsc, ColumnName: "first_name"},
				{Direction: basedto.DirectionAsc, ColumnName: "last_name"},
			},
		},
	}
}

// Validate implements interface basedto.ReqValidator
func (req *ListUserSimpleReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators,
		basedto.ValidateSlice(req.Status, true, 0, base.AllUserStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListUserSimpleResp struct {
	Meta *basedto.Meta     `json:"meta"`
	Data []*UserSimpleResp `json:"data"`
}

type UserSimpleResp struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Photo    string `json:"photo"`
}

func TransformUsersSimple(users []*entity.User) []*UserSimpleResp {
	return gofn.MapSlice(users, func(user *entity.User) *UserSimpleResp {
		return &UserSimpleResp{
			ID:       user.ID,
			FullName: user.FullName,
		}
	})
}
