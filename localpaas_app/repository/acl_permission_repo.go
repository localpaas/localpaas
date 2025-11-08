package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/tiendc/gofn"
	"github.com/uptrace/bun"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
)

type ACLPermissionRepo interface {
	ListByResources(ctx context.Context, db database.IDB, resources []*base.PermissionResource,
		opts ...bunex.SelectQueryOption) ([]*entity.ACLPermission, error)
	ListByUsers(ctx context.Context, db database.IDB, userIDs []string,
		opts ...bunex.SelectQueryOption) ([]*entity.ACLPermission, error)
	List(ctx context.Context, db database.IDB, paging *basedto.Paging,
		opts ...bunex.SelectQueryOption) ([]*entity.ACLPermission, *basedto.PagingMeta, error)

	UpsertMulti(ctx context.Context, db database.IDB, permissions []*entity.ACLPermission,
		conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error

	DeleteByResources(ctx context.Context, db database.IDB, resources []*base.PermissionResource,
		opts ...bunex.DeleteQueryOption) error
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

	// Construct the multi-column IN clause
	conditions := make([]string, 0, len(resources))
	args := make([]any, 0, len(resources)*2) //nolint:mnd
	for _, res := range resources {
		conditions = append(conditions, "(?,?)")
		args = append(args, res.SubjectID, res.ResourceID)
	}

	var permissions []*entity.ACLPermission
	query := db.NewSelect().Model(&permissions).
		Where(fmt.Sprintf("(subject_id,resource_id) IN (%s)", strings.Join(conditions, ",")), args...)
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
	query := db.NewSelect().Model(&permissions).Where("subject_id IN (?)", bun.In(userIDs))
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if err != nil {
		return nil, apperrors.New(err)
	}
	return permissions, nil
}

func (repo *aclPermissionRepo) List(ctx context.Context, db database.IDB, paging *basedto.Paging,
	opts ...bunex.SelectQueryOption) ([]*entity.ACLPermission, *basedto.PagingMeta, error) {
	var acls []*entity.ACLPermission
	query := db.NewSelect().Model(&acls)
	query = bunex.ApplySelect(query, opts...)

	pagingMeta := newPagingMeta(paging)

	// Counts the total first
	if paging != nil {
		total, err := query.Count(ctx)
		if err != nil {
			return nil, nil, apperrors.Wrap(err)
		}
		pagingMeta.Total = total
	}

	// Applies pagination
	query = bunex.ApplyPagination(query, paging)
	err := query.Scan(ctx)
	if err != nil {
		return nil, nil, wrapPaginationError(err, paging)
	}

	return acls, pagingMeta, nil
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

func (repo *aclPermissionRepo) DeleteByResources(ctx context.Context, db database.IDB,
	resources []*base.PermissionResource, opts ...bunex.DeleteQueryOption) error {
	if len(resources) == 0 {
		return nil
	}

	// Construct the multi-column IN clause
	conditions := make([]string, 0, len(resources))
	args := make([]any, 0, len(resources)*2) //nolint:mnd
	conditions2 := make([]string, 0, len(resources))
	args2 := make([]any, 0, len(resources)*2) //nolint:mnd
	for _, res := range resources {
		if res.ResourceID != "" {
			// Delete a specific row by a pair of (subject_id, resource_id)
			conditions = append(conditions, "(?,?)")
			args = append(args, res.SubjectID, res.ResourceID)
		} else if res.ResourceType != "" {
			// Delete multiple rows by a pair of (subject_id, resource_type)
			conditions2 = append(conditions2, "(?,?)")
			args2 = append(args2, res.SubjectID, res.ResourceType)
		}
	}

	query := db.NewDelete().Model((*entity.ACLPermission)(nil))
	if len(conditions) > 0 {
		query = query.Where(
			fmt.Sprintf("(subject_id,resource_id) IN (%s)", strings.Join(conditions, ",")), args...)
	}
	if len(conditions2) > 0 {
		where := gofn.If(len(conditions) > 0, query.WhereOr, query.Where) //nolint:staticcheck
		query = where(
			fmt.Sprintf("(subject_id,resource_type) IN (%s)", strings.Join(conditions2, ",")), args2...)
	}

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
		Where("subject_id IN (?)", bun.In(userIDs))
	query = bunex.ApplyDelete(query, opts...)

	_, err := query.Exec(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
