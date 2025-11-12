package useruc

import (
	"context"
	"encoding/base64"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity/appentity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/jwtsession"
	"github.com/localpaas/localpaas/localpaas_app/pkg/totp"
	"github.com/localpaas/localpaas/localpaas_app/usecase/useruc/userdto"
)

func (uc *UserUC) BeginUserSignup(
	ctx context.Context,
	req *userdto.BeginUserSignupReq,
) (*userdto.BeginUserSignupResp, error) {
	inviteToken := &appentity.UserInviteTokenClaims{}
	err := jwtsession.ParseToken(req.InviteToken, inviteToken)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrTokenInvalid).WithCause(err)
	}

	user, err := uc.userRepo.GetByID(ctx, uc.db, inviteToken.UserID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	if user.Status != base.UserStatusPending {
		return nil, apperrors.New(apperrors.ErrActionNotAllowed).
			WithMsgLog("user '%s' not require signup", user.Email)
	}

	resp := &userdto.BeginUserSignupDataResp{
		Username:       user.Username,
		Email:          user.Email,
		Role:           user.Role,
		SecurityOption: user.SecurityOption,
	}
	if !user.AccessExpireAt.IsZero() {
		resp.AccessExpiration = &user.AccessExpireAt
	}

	// Generate TOTP secret and QR code for user to setup 2FA
	if user.SecurityOption == base.UserSecurityPassword2FA {
		secret, qrCode, err := totp.GenerateSecretAndQRCode(qrCodeImageSize)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp.MFATotpSecret = secret
		resp.QRCode = &userdto.MFATotpQRCodeResp{
			DataBase64: base64.StdEncoding.EncodeToString(qrCode.Bytes()),
			ImageType:  qrCodeImageType,
			ImageSize:  qrCodeImageSize,
		}
	}

	return &userdto.BeginUserSignupResp{
		Data: resp,
	}, nil
}
