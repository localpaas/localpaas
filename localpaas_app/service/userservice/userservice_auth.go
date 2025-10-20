package userservice

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/jwtsession"
)

const (
	MFATokenExpiration = 60 * time.Second
)

type MFATokenClaims struct {
	jwtsession.BaseClaims
	UserID          string       `json:"userId"`
	MFAType         base.MFAType `json:"mfaType"`
	TrustedDeviceID string       `json:"deviceId,omitempty"`
}

// GenerateMFAToken builds MFA token for using in the next step
func (s *userService) GenerateMFAToken(userID string, mfaType base.MFAType, trustedDeviceID string) (string, error) {
	mfaToken, err := jwtsession.GenerateToken(&MFATokenClaims{
		UserID:          userID,
		MFAType:         mfaType,
		TrustedDeviceID: trustedDeviceID,
	}, MFATokenExpiration)
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	return mfaToken, nil
}
