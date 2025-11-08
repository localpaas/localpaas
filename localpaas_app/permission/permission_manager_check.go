package permission

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
)

type AccessCheck struct {
	RequireAdmin     bool
	SubjectType      base.SubjectType
	SubjectID        string
	ResourceType     base.ResourceType
	ResourceID       string
	ParentResourceID string
	Action           base.ActionType
}

func (p *manager) CheckAccess(ctx context.Context, db database.IDB, check *AccessCheck) (bool, error) {
	switch check.ResourceType { //nolint:exhaustive
	case base.ResourceTypeProject:
		return p.CheckProjectAccess(ctx, db, check)
	case base.ResourceTypeApp:
		return p.CheckAppAccess(ctx, db, check)
	}

	perms, err := p.aclPermissionRepo.ListByResources(ctx, db, []*base.PermissionResource{
		{
			SubjectType:  check.SubjectType,
			SubjectID:    check.SubjectID,
			ResourceType: check.ResourceType,
			ResourceID:   check.ResourceID,
		},
	})
	if err != nil || len(perms) == 0 {
		return false, apperrors.Wrap(err)
	}

	for _, perm := range perms {
		if check.Action == base.ActionTypeRead && perm.Actions.Read {
			return true, nil
		}
		if check.Action == base.ActionTypeWrite && perm.Actions.Write {
			return true, nil
		}
		if check.Action == base.ActionTypeDelete && perm.Actions.Delete {
			return true, nil
		}
	}

	return false, nil
}

func (p *manager) CheckProjectAccess(ctx context.Context, db database.IDB, check *AccessCheck) (bool, error) {
	acls, err := p.LoadProjectAccesses(ctx, db, check.ResourceID,
		bunex.SelectWhere("\"user\".id = ?", check.SubjectID),
	)
	if err != nil {
		return false, apperrors.Wrap(err)
	}
	return len(acls) > 0, nil
}

func (p *manager) CheckAppAccess(ctx context.Context, db database.IDB, check *AccessCheck) (bool, error) {
	acls, err := p.LoadAppAccesses(ctx, db, check.ParentResourceID, check.ResourceID,
		bunex.SelectWhere("\"user\".id = ?", check.SubjectID),
	)
	if err != nil {
		return false, apperrors.Wrap(err)
	}
	return len(acls) > 0, nil
}

func (p *manager) LoadProjectAccesses(ctx context.Context, db database.IDB, projectID string,
	extraLoadOpts ...bunex.SelectQueryOption) ([]*entity.ACLPermission, error) {
	loadOpts := []bunex.SelectQueryOption{
		bunex.SelectDistinct(),
		bunex.SelectWhere("\"user\".deleted_at IS NULL"),
		bunex.SelectWhere("\"user\".status = ?", base.UserStatusActive),
		bunex.SelectWhere("(\"user\".access_expire_at IS NULL OR \"user\".access_expire_at > NOW())"),

		bunex.SelectJoin("LEFT JOIN acl_permissions AS acl ON \"user\".id = acl.subject_id AND "+
			"acl.resource_id = ?", projectID),
		bunex.SelectWhereGroup(
			bunex.SelectWhere("\"user\".role = ?", base.UserRoleAdmin),
			bunex.SelectWhereOr("(acl.action_read OR acl.action_write OR acl.action_delete)"),
		),

		bunex.SelectRelation("Accesses",
			bunex.SelectWhere("acl_permission.deleted_at IS NULL"),
			bunex.SelectWhere("acl_permission.resource_id = ?", projectID),
		),
	}
	loadOpts = append(loadOpts, extraLoadOpts...)
	users, _, err := p.userRepo.List(ctx, db, nil, loadOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	aclPermissions := make([]*entity.ACLPermission, 0, len(users))
	for _, user := range users {
		var aclPerm *entity.ACLPermission
		if len(user.Accesses) > 0 {
			aclPerm = user.Accesses[0]
		}
		if user.Role == base.UserRoleAdmin {
			if aclPerm == nil {
				aclPerm = &entity.ACLPermission{
					SubjectType:  base.SubjectTypeUser,
					SubjectID:    user.ID,
					ResourceType: base.ResourceTypeProject,
					ResourceID:   projectID,
				}
			}
			aclPerm.Actions.Read = true
			aclPerm.Actions.Write = true
			aclPerm.Actions.Delete = true
		}
		if aclPerm != nil {
			aclPerm.SubjectUser = user
			aclPermissions = append(aclPermissions, aclPerm)
		}
	}

	return aclPermissions, nil
}

func (p *manager) LoadAppAccesses(ctx context.Context, db database.IDB, projectID, appID string,
	extraLoadOpts ...bunex.SelectQueryOption) ([]*entity.ACLPermission, error) {
	deniedACLs, _, err := p.aclPermissionRepo.List(ctx, db, nil,
		bunex.SelectWhere("deleted_at IS NULL"),
		bunex.SelectWhere("subject_type = ?", base.SubjectTypeUser),
		bunex.SelectWhere("resource_type = ?", base.ResourceTypeApp),
		bunex.SelectWhere("resource_id = ?", appID),
		bunex.SelectWhere("(action_read = FALSE AND action_write = FALSE AND action_delete = FALSE)"),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	deniedUserIDs := make([]string, 0, len(deniedACLs))
	for _, deniedACL := range deniedACLs {
		deniedUserIDs = append(deniedUserIDs, deniedACL.SubjectID)
	}
	if len(deniedUserIDs) == 0 {
		deniedUserIDs = []string{""}
	}

	loadOpts := []bunex.SelectQueryOption{
		bunex.SelectDistinct(),
		bunex.SelectWhere("\"user\".deleted_at IS NULL"),
		bunex.SelectWhere("\"user\".status = ?", base.UserStatusActive),
		bunex.SelectWhere("(\"user\".access_expire_at IS NULL OR \"user\".access_expire_at > NOW())"),

		bunex.SelectJoin("LEFT JOIN acl_permissions AS acl ON \"user\".id = acl.subject_id"),
		bunex.SelectWhereGroup(
			bunex.SelectWhere("\"user\".role = ?", base.UserRoleAdmin),
			// Has permission on the app
			bunex.SelectWhereOr("(acl.resource_id = ? AND "+
				"(acl.action_read OR acl.action_write OR acl.action_delete))", appID),
			// Has permission on the belonging project but not denied by the app
			bunex.SelectWhereOr("(acl.resource_id = ? AND "+
				"(acl.action_read OR acl.action_write OR acl.action_delete) AND acl.subject_id NOT IN (?))",
				projectID, bunex.In(deniedUserIDs)),
		),

		bunex.SelectRelation("Accesses",
			bunex.SelectWhere("acl_permission.deleted_at IS NULL"),
			bunex.SelectWhere("acl_permission.resource_id IN (?)", bunex.In([]string{projectID, appID})),
		),
	}
	loadOpts = append(loadOpts, extraLoadOpts...)
	users, _, err := p.userRepo.List(ctx, db, nil, loadOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	aclPermissions := make([]*entity.ACLPermission, 0, len(users))
	for _, user := range users {
		var aclPerm *entity.ACLPermission
		for _, acl := range user.Accesses {
			if acl.ResourceID == appID {
				aclPerm = acl
				break
			}
		}
		if aclPerm == nil && len(user.Accesses) > 0 {
			aclPerm = user.Accesses[0]
		}
		if user.Role == base.UserRoleAdmin {
			if aclPerm == nil {
				aclPerm = &entity.ACLPermission{
					SubjectType:  base.SubjectTypeUser,
					SubjectID:    user.ID,
					ResourceType: base.ResourceTypeApp,
					ResourceID:   appID,
				}
			}
			aclPerm.Actions.Read = true
			aclPerm.Actions.Write = true
			aclPerm.Actions.Delete = true
		}
		if aclPerm != nil {
			aclPerm.SubjectUser = user
			aclPermissions = append(aclPermissions, aclPerm)
		}
	}

	return aclPermissions, nil
}
