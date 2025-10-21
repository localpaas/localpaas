package permission

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type AccessCheck struct {
	UserID       string
	ResourceType base.ResourceType
	ResourceID   string
	Action       base.ActionType
}

func (p *manager) CheckAccess(ctx context.Context, db database.IDB, check *AccessCheck) (bool, error) {
	perms, err := p.aclPermissionRepo.ListByResources(ctx, db, []*base.PermissionResource{
		{
			UserID:       check.UserID,
			ResourceType: check.ResourceType,
			ResourceID:   check.ResourceID,
		},
	})
	if err != nil || len(perms) == 0 {
		return false, apperrors.Wrap(err)
	}

	for _, perm := range perms {
		if check.Action == base.ActionTypeRead && perm.Actions.Read == base.AccessTypeYes {
			return true, nil
		}
		if check.Action == base.ActionTypeWrite && perm.Actions.Write == base.AccessTypeYes {
			return true, nil
		}
		if check.Action == base.ActionTypeDelete && perm.Actions.Delete == base.AccessTypeYes {
			return true, nil
		}
	}

	return false, nil
}
