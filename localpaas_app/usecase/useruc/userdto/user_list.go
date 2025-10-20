package userdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type ListUserReq struct {
	Status []base.UserStatus `json:"-" mapstructure:"status"`
	Search string            `json:"-" mapstructure:"search"`

	Paging basedto.Paging `json:"-"`
}

func NewListUserReq() *ListUserReq {
	return &ListUserReq{
		Paging: basedto.Paging{
			// Default paging if unset by client
			Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "created_at"}},
		},
	}
}

func (req *ListUserReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators,
		basedto.ValidateSlice(req.Status, true, 0, base.AllUserStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListUserResp struct {
	Meta *basedto.Meta `json:"meta"`
	Data []*UserResp   `json:"data"`
}

func TransformUser(user *entity.User) (*UserResp, error) {
	resp := &UserResp{
		ID:             user.ID,
		Email:          user.Email,
		Role:           user.Role,
		Status:         user.Status,
		FullName:       user.FullName,
		Photo:          user.Photo,
		SecurityOption: user.SecurityOption,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
	if !user.AccessExpireAt.IsZero() {
		resp.AccessExpireAt = &user.AccessExpireAt
	}
	if !user.LastAccess.IsZero() {
		resp.LastAccess = &user.LastAccess
	}
	return resp, nil
}

func TransformUsers(users []*entity.User) ([]*UserResp, error) {
	return basedto.TransformObjectSlice(users, TransformUser) //nolint:wrapcheck
}
