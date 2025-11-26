package useruc

import (
	"context"
	"fmt"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/usecase/useruc/userdto"
)

const (
	//nolint passwordResetFEPath front-end path to the UI
	passwordResetFEPath = "auth/reset-password"
)

func (uc *UserUC) RequestResetPassword(
	ctx context.Context,
	auth *basedto.Auth,
	req *userdto.RequestResetPasswordReq,
) (*userdto.RequestResetPasswordResp, error) {
	user, err := uc.userRepo.GetByID(ctx, uc.db, auth.User.ID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	if user.SecurityOption == base.UserSecurityEnforceSSO {
		return nil, apperrors.New(apperrors.ErrActionNotAllowed).
			WithMsgLog("user authentication method is enforce-sso")
	}

	token, err := uc.userService.GeneratePasswordResetToken(user.ID)
	if err != nil {
		return nil, apperrors.New(err).WithMsgLog("failed to generate password reset token")
	}

	resetLink := fmt.Sprintf("%s/%s?userId=%s&token=%s", config.Current.BaseURL,
		passwordResetFEPath, user.ID, token)

	// TODO: handle req.SendEmail when email account is setup

	return &userdto.RequestResetPasswordResp{
		Data: &userdto.RequestResetPasswordDataResp{
			ResetPasswordLink: resetLink,
		},
	}, nil
}
