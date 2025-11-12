package sessionuc

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sessionuc/sessiondto"
)

const (
	// We allow at most 5 attempts of login in the first minute
	// The duration increase by exponential of 2 after each minute
	maxPasswordFailsInARow       = 5
	passwordCheckDurationEachRow = time.Minute
)

const (
	nextStepMfaInput = "NextMfa"
	nextStepMfaSetup = "NextMfaSetup"
)

func (uc *SessionUC) LoginWithPassword(
	ctx context.Context,
	req *sessiondto.LoginWithPasswordReq,
) (resp *sessiondto.LoginWithPasswordResp, err error) {
	dbUser, err := uc.userRepo.GetByUsernameOrEmail(ctx, uc.db, req.Username, req.Username)
	if err != nil {
		return nil, uc.wrapSensitiveError(err)
	}

	if err = uc.allowPasswordLoginAtTheMoment(dbUser); err != nil {
		return nil, err
	}

	err = uc.userService.VerifyPassword(dbUser, req.Password)
	_ = uc.savePasswordCheckingStatus(ctx, dbUser, err == nil)
	if err != nil {
		return nil, uc.wrapSensitiveError(err)
	}

	passcodeRequired := dbUser.TotpSecret != ""

	// When trusted device is sent
	if passcodeRequired && req.TrustedDeviceID != "" {
		timeNow := timeutil.NowUTC()
		// If the sending trusted device matches the data in DB
		trustedDevice, err := uc.loginTrustedDeviceRepo.GetByUserAndDevice(ctx, uc.db, dbUser.ID, req.TrustedDeviceID)
		if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.Wrap(err)
		}
		if trustedDevice != nil && timeNow.Sub(trustedDevice.UpdatedAt) < config.Current.Session.DeviceTrustedPeriod {
			passcodeRequired = false
		}
	}

	// When passcode is required, builds token for using in the next step
	if passcodeRequired {
		mfaType := base.MFATypeTOTP
		mfaToken, err := uc.userService.GenerateMFAToken(dbUser.ID, mfaType, req.TrustedDeviceID)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}

		return &sessiondto.LoginWithPasswordResp{
			Data: &sessiondto.LoginWithPasswordDataResp{
				NextStep: nextStepMfaInput,
				MFAType:  mfaType,
				MFAToken: mfaToken,
			},
		}, nil
	}

	// Create a new session as login succeeds
	sessionData, err := uc.createSession(ctx, &sessiondto.BaseCreateSessionReq{User: dbUser})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	var nextStep string
	if dbUser.SecurityOption == base.UserSecurityPassword2FA && dbUser.TotpSecret == "" {
		nextStep = nextStepMfaSetup
	}

	return &sessiondto.LoginWithPasswordResp{
		Data: &sessiondto.LoginWithPasswordDataResp{
			Session:  sessionData,
			NextStep: nextStep,
		},
	}, nil
}

// allowPasswordLoginAtTheMoment checks if user can do password login at the moment.
// If user made too many login failures, they need to wait for some time before they can try again.
func (uc *SessionUC) allowPasswordLoginAtTheMoment(dbUser *entity.User) error {
	if dbUser.SecurityOption == base.UserSecurityEnforceSSO {
		return apperrors.New(apperrors.ErrSSORequired)
	}
	if dbUser.PasswordFailsInRow < maxPasswordFailsInARow {
		return nil
	}
	expo := dbUser.PasswordFailsInRow / maxPasswordFailsInARow
	minWaitingDuration := time.Duration(math.Pow(2, float64(expo))) * passwordCheckDurationEachRow //nolint:mnd
	durationFromFirstFail := timeutil.NowUTC().Sub(dbUser.PasswordFirstFailAt)
	if durationFromFirstFail > minWaitingDuration {
		return nil
	}
	waitingDuration := int((minWaitingDuration - durationFromFirstFail).Seconds())
	return apperrors.New(apperrors.ErrTooManyLoginFailures).WithParam("WaitDuration", waitingDuration)
}

// savePasswordCheckingStatus saves password checking status including the number of failures
// and timestamp of the first fail.
func (uc *SessionUC) savePasswordCheckingStatus(ctx context.Context, dbUser *entity.User, success bool) error {
	passwordFailsInRow := dbUser.PasswordFailsInRow
	passwordFirstFailAt := dbUser.PasswordFirstFailAt
	if success {
		// If user has password check fails, clear it
		passwordFailsInRow = 0
		passwordFirstFailAt = time.Time{}
	} else {
		// Save failed check count and update the first fail timestamp
		passwordFailsInRow++
		if passwordFirstFailAt.IsZero() {
			passwordFirstFailAt = time.Now()
		}
	}
	// Updates user data into DB if there is any change
	if passwordFailsInRow == dbUser.PasswordFailsInRow && passwordFirstFailAt.Equal(dbUser.PasswordFirstFailAt) {
		return nil
	}
	dbUser.PasswordFailsInRow = passwordFailsInRow
	dbUser.PasswordFirstFailAt = passwordFirstFailAt
	err := uc.userRepo.Update(ctx, uc.db, dbUser,
		bunex.UpdateColumns("password_fails_in_row", "password_first_fail_at"))
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (uc *SessionUC) wrapSensitiveError(err error) error {
	// Due to security reason, we don't want to send the real error to user for the cases
	// user not found and password mismatched.
	if errors.Is(err, apperrors.ErrNotFound) || errors.Is(err, apperrors.ErrPasswordMismatched) ||
		errors.Is(err, apperrors.ErrAPIKeyMismatched) {
		// Notes that the `cause` only shows up in dev env, not in production
		return apperrors.New(apperrors.ErrLoginInputInvalid).WithCause(err)
	}
	return apperrors.Wrap(err)
}
