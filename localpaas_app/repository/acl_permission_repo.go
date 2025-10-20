package repository

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
)

type ACLPermissionRepo interface {
	ListByResources(ctx context.Context, db database.IDB, resources []*base.PermissionResource,
		opts ...bunex.SelectQueryOption) ([]*entity.ACLPermission, error)
	ListByUsers(ctx context.Context, db database.IDB, userIDs []string,
		opts ...bunex.SelectQueryOption) ([]*entity.ACLPermission, error)

	UpsertMulti(ctx context.Context, db database.IDB, permissions []*entity.ACLPermission,
		conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error

	DeleteByIDs(ctx context.Context, db database.IDB, ids []string, opts ...bunex.DeleteQueryOption) error
	DeleteByUsers(ctx context.Context, db database.IDB, userIDs []string,
		opts ...bunex.DeleteQueryOption) error
}

type aclPermissionRepo struct {
}

func NewACLPermissionRepo() ACLPermissionRepo {
	return &aclPermissionRepo{}
}

func (repo *aclPermissionRepo) ListByResources(ctx context.Context, db database.IDB,
	resources []*base.PermissionResource, opts ...bunex.SelectQueryOption) ([]*entity.ACLPermission, error) {
	if len(resources) == 0 {
		return nil, nil
	}

	// opts = append(opts, bunex.SelectWhereGroup(
	//	lo.Map(resources, func(res *base.PermissionResource, index int) bunex.SelectQueryOption {
	//		if index == 0 {
	//			return bunex.SelectWhere("(user_id,resource_type,resource_id) = (?,?,?)",
	//				res.UserID, res.ResourceType, res.ResourceID)
	//		}
	//		return bunex.SelectWhereOr("(user_id,resource_type,resource_id) = (?,?,?)",
	//			res.UserID, res.ResourceType, res.ResourceID)
	//	})...,
	// ))

	// Construct the multi-column IN clause
	conditions := make([]string, 0, len(resources))
	args := make([]any, 0, len(resources)*3) //nolint:mnd
	for _, res := range resources {
		conditions = append(conditions, "(?,?,?)")
		args = append(args, res.UserID, res.ResourceType, res.ResourceID)
	}

	var permissions []*entity.ACLPermission
	query := db.NewSelect().Model(&permissions).
		Where(fmt.Sprintf("(user_id,resource_type,resource_id) IN (%s)", bun.In(conditions)), args...)
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if err != nil {
		return nil, apperrors.New(err)
	}
	return permissions, nil
}

func (repo *aclPermissionRepo) ListByUsers(ctx context.Context, db database.IDB, userIDs []string,
	opts ...bunex.SelectQueryOption) ([]*entity.ACLPermission, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	var permissions []*entity.ACLPermission
	query := db.NewSelect().Model(&permissions).Where("user_id IN (?)", bun.In(userIDs))
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if err != nil {
		return nil, apperrors.New(err)
	}
	return permissions, nil
}

func (repo *aclPermissionRepo) UpsertMulti(ctx context.Context, db database.IDB, permissions []*entity.ACLPermission,
	conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error {
	if len(permissions) == 0 {
		return nil
	}
	query := db.NewInsert().Model(&permissions)
	query = bunex.ApplyInsert(query, opts...)
	query = bunex.ApplyUpsert(query, conflictCols, updateCols)

	_, err := query.Exec(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (repo *aclPermissionRepo) DeleteByIDs(ctx context.Context, db database.IDB, ids []string,
	opts ...bunex.DeleteQueryOption) error {
	if len(ids) == 0 {
		return nil
	}
	query := db.NewDelete().Model((*entity.ACLPermission)(nil)).
		Where("id IN (?)", bun.In(ids))
	query = bunex.ApplyDelete(query, opts...)

	_, err := query.Exec(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (repo *aclPermissionRepo) DeleteByUsers(ctx context.Context, db database.IDB, userIDs []string,
	opts ...bunex.DeleteQueryOption) error {
	if len(userIDs) == 0 {
		return nil
	}
	query := db.NewDelete().Model((*entity.ACLPermission)(nil)).
		Where("user_id IN (?)", bun.In(userIDs))
	query = bunex.ApplyDelete(query, opts...)

	_, err := query.Exec(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
