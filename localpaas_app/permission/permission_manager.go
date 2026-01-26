package permission

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/repository"
)

type Manager interface {
	CheckAccess(ctx context.Context, db database.IDB, auth *basedto.Auth, check *AccessCheck) (bool, error)

	// NOTE: this func should be called within a transaction
	UpdateACLPermissions(ctx context.Context, db database.IDB, perms []*entity.ACLPermission) error
	RemoveACLPermissions(ctx context.Context, db database.IDB, perms []*base.PermissionResource) error

	LoadObjectAccesses(ctx context.Context, db database.IDB, check *AccessCheck, sort bool,
		extraLoadOpts ...bunex.SelectQueryOption) ([]*entity.ACLPermission, error)
}

type manager struct {
	aclPermissionRepo repository.ACLPermissionRepo
	userRepo          repository.UserRepo
	projectRepo       repository.ProjectRepo
}

func NewManager(
	aclPermissionRepo repository.ACLPermissionRepo,
	userRepo repository.UserRepo,
	projectRepo repository.ProjectRepo,
) Manager {
	return &manager{
		aclPermissionRepo: aclPermissionRepo,
		userRepo:          userRepo,
		projectRepo:       projectRepo,
	}
}
