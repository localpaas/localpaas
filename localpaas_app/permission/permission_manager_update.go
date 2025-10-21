package permission

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

func (p *manager) UpdateACLPermissions(ctx context.Context, db database.IDB, perms []*entity.ACLPermission) error {
	var toUpsert []*entity.ACLPermission
	var toDelete []*base.PermissionResource
	for _, perm := range perms {
		if perm.Actions.IsNoAccess() {
			toDelete = append(toDelete, &base.PermissionResource{
				UserID:       perm.UserID,
				ResourceType: perm.ResourceType,
				ResourceID:   perm.ResourceID,
			})
			continue
		}
		toUpsert = append(toUpsert, perm)
	}

	err := p.aclPermissionRepo.DeleteByResources(ctx, db, toDelete)
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = p.aclPermissionRepo.UpsertMulti(ctx, db, toUpsert,
		entity.ACLPermissionUpsertingConflictCols, entity.ACLPermissionUpsertingUpdateCols)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
