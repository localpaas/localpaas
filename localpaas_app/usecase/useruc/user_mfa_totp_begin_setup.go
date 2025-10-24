package useruc

import (
	"context"
	"encoding/base64"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/useruc/userdto"
	"github.com/localpaas/localpaas/pkg/totp"
)

const (
	qrCodeImageSize = 300
	qrCodeImageType = "image/png"
)

func (uc *UserUC) BeginMFATotpSetup(
	ctx context.Context,
	auth *basedto.Auth,
	req *userdto.BeginMFATotpSetupReq,
) (*userdto.BeginMFATotpSetupResp, error) {
	user, err := uc.userRepo.GetByID(ctx, uc.db, auth.User.ID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	if user.SecurityOption == base.UserSecurityEnforceSSO {
		return nil, apperrors.New(apperrors.ErrActionNotAllowed).
			WithMsgLog("user authentication method is enforce-sso")
	}

	secret, qrCode, err := totp.GenerateSecretAndQRCode(qrCodeImageSize)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	totpToken, err := uc.userService.GenerateMFATotpSetupToken(user.ID, secret)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &userdto.BeginMFATotpSetupResp{
		Data: &userdto.MFATotpSetupDataResp{
			Secret:    secret,
			TotpToken: totpToken,
			QRCode: &userdto.MFATotpSetupQRCodeResp{
				DataBase64: base64.StdEncoding.EncodeToString(qrCode.Bytes()),
				ImageType:  qrCodeImageType,
				ImageSize:  qrCodeImageSize,
			},
		},
	}, nil
}
