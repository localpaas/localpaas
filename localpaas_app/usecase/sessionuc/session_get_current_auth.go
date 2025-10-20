package sessionuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/pkg/timeutil"
)

func (uc *SessionUC) GetCurrentAuth(ctx context.Context, jwt string) (*basedto.Auth, error) {
	user, err := uc.GetCurrentUser(ctx, jwt)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	auth := &basedto.Auth{User: user}

	// User status is not active
	if user.Status != base.UserStatusActive {
		return nil, apperrors.New(apperrors.ErrUserUnavailable).
			WithMsgLog("user status: %s", user.Status)
	}

	// User must have access permission
	if user.IsAccessExpired() {
		return nil, apperrors.New(apperrors.ErrUserUnavailable).
			WithMsgLog("user access expired at: %v", user.AccessExpireAt)
	}

	// Update `last_access` timestamp after each period of few minute.
	// NOTE: We can't update `last_access` timestamp every request due to performance reason.
	timeNow := timeutil.NowUTC()
	if user.LastAccess.IsZero() ||
		timeNow.Sub(user.LastAccess) > config.Current.Session.LastAccessUpdatePeriod {
		user.LastAccess = timeNow
		// Just ignore the error if happens (as this is not important action)
		_ = uc.userRepo.Update(ctx, uc.db, user.Entity(), bunex.UpdateColumns("last_access"))
	}

	return auth, nil
}
