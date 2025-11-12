package sessiondto

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
)

type UserResp struct {
	ID             string                  `json:"id"`
	Username       string                  `json:"username"`
	Email          string                  `json:"email"`
	Role           base.UserRole           `json:"role"`
	Status         base.UserStatus         `json:"status"`
	FullName       string                  `json:"fullName"`
	Photo          string                  `json:"photo"`
	SecurityOption base.UserSecurityOption `json:"securityOption"`
	MfaSecret      string                  `json:"mfaSecret"`
	CreatedAt      time.Time               `json:"createdAt"`
	UpdatedAt      time.Time               `json:"updatedAt"`
	AccessExpireAt *timeutil.Date          `json:"accessExpireAt" copy:",nilonzero"`
}

func TransformUser(user *entity.User) (resp *UserResp, err error) {
	if err = copier.Copy(&resp, &user); err != nil {
		return nil, apperrors.Wrap(err)
	}
	if user.TotpSecret != "" {
		resp.MfaSecret = "********"
	}
	return resp, nil
}
