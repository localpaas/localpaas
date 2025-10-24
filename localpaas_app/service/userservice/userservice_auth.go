package userservice

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity/appentity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/jwtsession"
)

const (
	MFATokenExpiration          = 60 * time.Second
	MFATotpSetupTokenExpiration = 120 * time.Second
)

// GenerateMFAToken builds MFA token for using in the next step
func (s *userService) GenerateMFAToken(userID string, mfaType base.MFAType, trustedDeviceID string) (string, error) {
	mfaToken, err := jwtsession.GenerateToken(&appentity.MFATokenClaims{
		UserID:          userID,
		MFAType:         mfaType,
		TrustedDeviceID: trustedDeviceID,
	}, MFATokenExpiration)
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	return mfaToken, nil
}

// GenerateMFATotpSetupToken builds MFA TOTP token for setting up TOTP secret
func (s *userService) GenerateMFATotpSetupToken(userID string, toptSecret string) (string, error) {
	mfaToken, err := jwtsession.GenerateToken(&appentity.MFATotpSetupTokenClaims{
		UserID: userID,
		Secret: toptSecret,
	}, MFATotpSetupTokenExpiration)
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	return mfaToken, nil
}
