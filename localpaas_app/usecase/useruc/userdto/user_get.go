package userdto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
)

type GetUserReq struct {
	ID string `json:"-"`
}

func NewGetUserReq() *GetUserReq {
	return &GetUserReq{}
}

func (req *GetUserReq) Validate() apperrors.ValidationErrors {
	return apperrors.NewValidationErrors(vld.Validate(basedto.ValidateID(&req.ID, true, "id")...))
}

type GetUserResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *UserDetailsResp  `json:"data"`
}

type UserDetailsResp struct {
	*UserResp
}

type UserResp struct {
	ID             string                  `json:"id"`
	Username       string                  `json:"username"`
	Email          string                  `json:"email"`
	Role           base.UserRole           `json:"role"`
	Status         base.UserStatus         `json:"status"`
	FullName       string                  `json:"fullName"`
	Photo          string                  `json:"photo"`
	Position       string                  `json:"position"`
	SecurityOption base.UserSecurityOption `json:"securityOption"`
	Notes          string                  `json:"notes,omitempty"`

	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	AccessExpireAt *time.Time `json:"accessExpireAt" copy:",nilonzero"`
	LastAccess     *time.Time `json:"lastAccess" copy:",nilonzero"`
}

func TransformUserDetails(user *entity.User) (resp *UserDetailsResp, err error) {
	var userResp *UserResp
	if err = copier.Copy(&userResp, &user); err != nil {
		return nil, apperrors.Wrap(err)
	}
	return &UserDetailsResp{
		UserResp: userResp,
	}, nil
}
