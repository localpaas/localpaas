package userservice

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity/appentity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/jwtsession"
)

const (
	MFATokenExp           = 60 * time.Second
	MFATotpSetupTokenExp  = 180 * time.Second
	UserInviteTokenExp    = 7 * 24 * time.Hour // 1 week
	PasswordResetTokenExp = 7 * 24 * time.Hour // 1 week
)

// GenerateMFAToken builds MFA token for using in the next step
func (s *userService) GenerateMFAToken(userID string, mfaType base.MFAType, trustedDeviceID string) (string, error) {
	mfaToken, err := jwtsession.GenerateToken(&appentity.MFATokenClaims{
		Kind:            "mfa",
		UserID:          userID,
		MFAType:         mfaType,
		TrustedDeviceID: trustedDeviceID,
	}, MFATokenExp)
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	return mfaToken, nil
}

func (s *userService) ParseMFAToken(token string) (*appentity.MFATokenClaims, error) {
	tokenClaims := &appentity.MFATokenClaims{}
	if err := jwtsession.ParseToken(token, tokenClaims); err != nil {
		return nil, apperrors.New(apperrors.ErrTokenInvalid).WithCause(err)
	}
	if tokenClaims.Kind != "mfa" {
		return nil, apperrors.New(apperrors.ErrTokenInvalid)
	}
	return tokenClaims, nil
}

// GenerateMFATotpSetupToken builds MFA TOTP token for setting up TOTP secret
func (s *userService) GenerateMFATotpSetupToken(userID string, toptSecret string) (string, error) {
	mfaToken, err := jwtsession.GenerateToken(&appentity.MFATotpSetupTokenClaims{
		Kind:   "mfa-setup",
		UserID: userID,
		Secret: toptSecret,
	}, MFATotpSetupTokenExp)
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	return mfaToken, nil
}

func (s *userService) ParseMFATotpSetupToken(token string) (*appentity.MFATotpSetupTokenClaims, error) {
	tokenClaims := &appentity.MFATotpSetupTokenClaims{}
	if err := jwtsession.ParseToken(token, tokenClaims); err != nil {
		return nil, apperrors.New(apperrors.ErrTokenInvalid).WithCause(err)
	}
	if tokenClaims.Kind != "mfa-setup" {
		return nil, apperrors.New(apperrors.ErrTokenInvalid)
	}
	return tokenClaims, nil
}

// GenerateUserInviteToken builds token for inviting users
func (s *userService) GenerateUserInviteToken(userID string) (string, error) {
	token, err := jwtsession.GenerateToken(&appentity.UserInviteTokenClaims{
		Kind:   "user-invite",
		UserID: userID,
	}, UserInviteTokenExp)
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	return token, nil
}

func (s *userService) ParseUserInviteToken(token string) (*appentity.UserInviteTokenClaims, error) {
	tokenClaims := &appentity.UserInviteTokenClaims{}
	if err := jwtsession.ParseToken(token, tokenClaims); err != nil {
		return nil, apperrors.New(apperrors.ErrTokenInvalid).WithCause(err)
	}
	if tokenClaims.Kind != "user-invite" {
		return nil, apperrors.New(apperrors.ErrTokenInvalid)
	}
	return tokenClaims, nil
}

// GeneratePasswordResetToken builds token for resetting passwords
func (s *userService) GeneratePasswordResetToken(userID string) (string, error) {
	token, err := jwtsession.GenerateToken(&appentity.PasswordResetTokenClaims{
		Kind:   "pwd-reset",
		UserID: userID,
	}, PasswordResetTokenExp)
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	return token, nil
}

func (s *userService) ParsePasswordResetToken(token string) (*appentity.PasswordResetTokenClaims, error) {
	tokenClaims := &appentity.PasswordResetTokenClaims{}
	if err := jwtsession.ParseToken(token, tokenClaims); err != nil {
		return nil, apperrors.New(apperrors.ErrTokenInvalid).WithCause(err)
	}
	if tokenClaims.Kind != "pwd-reset" {
		return nil, apperrors.New(apperrors.ErrTokenInvalid)
	}
	return tokenClaims, nil
}
