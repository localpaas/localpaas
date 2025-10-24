package sessionuc

import (
	"context"
	"errors"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/entity/appentity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/jwtsession"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sessionuc/sessiondto"
	"github.com/localpaas/localpaas/pkg/timeutil"
	"github.com/localpaas/localpaas/pkg/totp"
)

const (
	passcodeMaxAttempts = 5
)

func (uc *SessionUC) LoginWithPasscode(
	ctx context.Context,
	req *sessiondto.LoginWithPasscodeReq,
) (resp *sessiondto.LoginWithPasscodeResp, err error) {
	mfaTokenClaims := &appentity.MFATokenClaims{}
	if err = jwtsession.ParseToken(req.MFAToken, mfaTokenClaims); err != nil {
		return nil, apperrors.New(apperrors.ErrMFATokenInvalid).WithCause(err)
	}

	dbUser, err := uc.userRepo.GetByID(ctx, uc.db, mfaTokenClaims.UserID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	// Verify passcode TOTP
	if !totp.VerifyPasscode(req.Passcode, dbUser.TotpSecret) {
		passcode, err := uc.mfaPasscodeRepo.Get(ctx, mfaTokenClaims.UserID)
		if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.Wrap(err)
		}
		if passcode != nil && passcode.Attempts >= passcodeMaxAttempts {
			_ = uc.mfaPasscodeRepo.Del(ctx, mfaTokenClaims.UserID)
			return nil, apperrors.New(apperrors.ErrTooManyPasscodeAttempts).
				WithMsgLog("too many passcode attempts: %d", passcode.Attempts)
		}
		// Increase the attempts
		_ = uc.mfaPasscodeRepo.IncrAttempts(ctx, mfaTokenClaims.UserID, passcode)
		return nil, apperrors.Wrap(apperrors.ErrPasscodeMismatched)
	}

	// Removes the passcode in redis
	if err = uc.mfaPasscodeRepo.Del(ctx, mfaTokenClaims.UserID); err != nil {
		return nil, apperrors.Wrap(err)
	}

	// Save trusted device if needs to
	if mfaTokenClaims.TrustedDeviceID != "" {
		timeNow := timeutil.NowUTC()
		trustedDevice := &entity.LoginTrustedDevice{
			UserID:    dbUser.ID,
			DeviceID:  mfaTokenClaims.TrustedDeviceID,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		}
		err = uc.loginTrustedDeviceRepo.Upsert(ctx, uc.db, trustedDevice,
			entity.LoginTrustedDeviceUpsertingConflictCols, entity.LoginTrustedDeviceUpsertingUpdateCols)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
	}

	// Create a new session as login succeeds
	sessionData, err := uc.createSession(ctx, &sessiondto.BaseCreateSessionReq{User: dbUser})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &sessiondto.LoginWithPasscodeResp{
		Data: &sessiondto.LoginWithPasscodeDataResp{
			Session: sessionData,
		},
	}, nil
}
