package sessionuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/permission"
)

func (uc *SessionUC) VerifyAuth(ctx context.Context, auth *basedto.Auth, accessCheck *permission.AccessCheck) error {
	if auth.User.Role == base.UserRoleOwner || auth.User.Role == base.UserRoleAdmin {
		return nil
	}
	if accessCheck != nil && accessCheck.UserID == "" {
		accessCheck.UserID = auth.User.ID
	}
	checkResult, err := uc.permissionManager.CheckAccess(ctx, uc.db, accessCheck)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if !checkResult {
		return apperrors.New(apperrors.ErrUnauthorized)
	}
	return nil
}
