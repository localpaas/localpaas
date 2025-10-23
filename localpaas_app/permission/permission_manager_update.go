package permission

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

func (p *manager) UpdateACLPermissions(ctx context.Context, db database.IDB, perms []*entity.ACLPermission) error {
	err := p.aclPermissionRepo.UpsertMulti(ctx, db, perms,
		entity.ACLPermissionUpsertingConflictCols, entity.ACLPermissionUpsertingUpdateCols)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
