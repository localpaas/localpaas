package repository

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
)

type ProjectEnvRepo interface {
	List(ctx context.Context, db database.IDB, paging *basedto.Paging,
		opts ...bunex.SelectQueryOption) ([]*entity.ProjectEnv, *basedto.PagingMeta, error)

	Upsert(ctx context.Context, db database.IDB, projectEnv *entity.ProjectEnv,
		conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error
	UpsertMulti(ctx context.Context, db database.IDB, projectEnvs []*entity.ProjectEnv,
		conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error

	DeleteAllByProjects(ctx context.Context, db database.IDB, projectIDs []string,
		opts ...bunex.DeleteQueryOption) error
}

type projectEnvRepo struct {
}

func NewProjectEnvRepo() ProjectEnvRepo {
	return &projectEnvRepo{}
}

func (repo *projectEnvRepo) List(ctx context.Context, db database.IDB, paging *basedto.Paging,
	opts ...bunex.SelectQueryOption) ([]*entity.ProjectEnv, *basedto.PagingMeta, error) {
	var projectEnvs []*entity.ProjectEnv
	query := db.NewSelect().Model(&projectEnvs)
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

	return projectEnvs, pagingMeta, nil
}

func (repo *projectEnvRepo) Upsert(ctx context.Context, db database.IDB, projectEnv *entity.ProjectEnv,
	conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error {
	return repo.UpsertMulti(ctx, db, []*entity.ProjectEnv{projectEnv}, conflictCols, updateCols, opts...)
}

func (repo *projectEnvRepo) UpsertMulti(ctx context.Context, db database.IDB, projectEnvs []*entity.ProjectEnv,
	conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error {
	if len(projectEnvs) == 0 {
		return nil
	}
	query := db.NewInsert().Model(&projectEnvs)
	query = bunex.ApplyInsert(query, opts...)
	query = bunex.ApplyUpsert(query, conflictCols, updateCols)

	_, err := query.Exec(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (repo *projectEnvRepo) DeleteAllByProjects(ctx context.Context, db database.IDB,
	projectIDs []string, opts ...bunex.DeleteQueryOption) error {
	if len(projectIDs) == 0 {
		return nil
	}
	query := db.NewDelete().Model((*entity.ProjectEnv)(nil)).
		Where("project_id IN (?)", bun.In(projectIDs))
	query = bunex.ApplyDelete(query, opts...)

	_, err := query.Exec(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
