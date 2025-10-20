package sessionuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/permission"
)

func (uc *SessionUC) VerifyAuth(ctx context.Context, auth *basedto.Auth, accessCheck *permission.AccessCheck) error {
	// checkResult, err := uc.permissionChecker.CheckAccess(ctx, uc.db, &permission.AccessCheck{
	//	UserAccessInfo: &permission.UserAccessInfo{
	//		WorkspaceUserID: auth.WorkspaceUser.ID,
	//		RoleIDs:         roleIDs,
	//	},
	//	ComponentType: accessCheck.ComponentType,
	//	ResourceType:  accessCheck.ResourceType,
	//	ResourceID:    accessCheck.ResourceID,
	//	Action:        accessCheck.Action,
	// })
	// if err != nil {
	//	return nil, apperrors.Wrap(err)
	// }
	//
	// return checkResult, nil
	return nil
}
