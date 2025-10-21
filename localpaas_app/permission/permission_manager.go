package permission

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
)

type Manager interface {
	CheckAccess(ctx context.Context, db database.IDB, check *AccessCheck) (bool, error)

	// NOTE: this func should be called within a transaction
	UpdateACLPermissions(ctx context.Context, db database.IDB, perms []*entity.ACLPermission) error
}

type manager struct {
	aclPermissionRepo repository.ACLPermissionRepo
}

func NewManager(
	aclPermissionRepo repository.ACLPermissionRepo,
) Manager {
	return &manager{
		aclPermissionRepo: aclPermissionRepo,
	}
}
