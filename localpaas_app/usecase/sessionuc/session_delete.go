package sessionuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sessionuc/sessiondto"
)

func (uc *SessionUC) DeleteSession(
	ctx context.Context,
	user *basedto.User,
) (resp *sessiondto.DeleteSessionResp, err error) {
	// Invalidate the old token to make it unusable
	err = uc.userTokenRepo.Del(ctx, user.AuthClaims.UserID, user.AuthClaims.UID)
	if err != nil {
		return nil, apperrors.New(err).WithMsgLog("failed to invalidate old token")
	}

	return &sessiondto.DeleteSessionResp{}, nil
}
