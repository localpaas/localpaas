package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/uptrace/bun"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
)

type S3StorageRepo interface {
	GetByID(ctx context.Context, db database.IDB, id string,
		opts ...bunex.SelectQueryOption) (*entity.S3Storage, error)
	GetByName(ctx context.Context, db database.IDB, name string,
		opts ...bunex.SelectQueryOption) (*entity.S3Storage, error)
	List(ctx context.Context, db database.IDB, paging *basedto.Paging,
		opts ...bunex.SelectQueryOption) ([]*entity.S3Storage, *basedto.PagingMeta, error)
	ListByIDs(ctx context.Context, db database.IDB, ids []string,
		opts ...bunex.SelectQueryOption) ([]*entity.S3Storage, error)

	Upsert(ctx context.Context, db database.IDB, s3Storage *entity.S3Storage,
		conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error
	UpsertMulti(ctx context.Context, db database.IDB, s3Storages []*entity.S3Storage,
		conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error
}

type s3StorageRepo struct {
}

func NewS3StorageRepo() S3StorageRepo {
	return &s3StorageRepo{}
}

func (repo *s3StorageRepo) GetByID(ctx context.Context, db database.IDB, id string,
	opts ...bunex.SelectQueryOption) (*entity.S3Storage, error) {
	s3Storage := &entity.S3Storage{}
	query := db.NewSelect().Model(s3Storage).Where("s3_storage.id = ?", id)
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if s3Storage == nil || errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.NewNotFound("S3Storage").WithCause(err)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return s3Storage, nil
}

func (repo *s3StorageRepo) GetByName(ctx context.Context, db database.IDB, name string,
	opts ...bunex.SelectQueryOption) (*entity.S3Storage, error) {
	if name == "" {
		return nil, nil
	}
	s3Storage := &entity.S3Storage{}
	query := db.NewSelect().Model(s3Storage).Where("LOWER(s3_storage.name) = ?", strings.ToLower(name))
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if s3Storage == nil || errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.NewNotFound("S3Storage").WithCause(err)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return s3Storage, nil
}

func (repo *s3StorageRepo) List(ctx context.Context, db database.IDB, paging *basedto.Paging,
	opts ...bunex.SelectQueryOption) ([]*entity.S3Storage, *basedto.PagingMeta, error) {
	var s3Storages []*entity.S3Storage
	query := db.NewSelect().Model(&s3Storages)
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

	return s3Storages, pagingMeta, nil
}

func (repo *s3StorageRepo) ListByIDs(ctx context.Context, db database.IDB, ids []string,
	opts ...bunex.SelectQueryOption) ([]*entity.S3Storage, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var s3Storages []*entity.S3Storage
	query := db.NewSelect().Model(&s3Storages).Where("s3_storage.id IN (?)", bun.In(ids))
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if err != nil {
		return nil, apperrors.New(err)
	}
	return s3Storages, nil
}

func (repo *s3StorageRepo) Upsert(ctx context.Context, db database.IDB, s3Storage *entity.S3Storage,
	conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error {
	return repo.UpsertMulti(ctx, db, []*entity.S3Storage{s3Storage}, conflictCols, updateCols, opts...)
}

func (repo *s3StorageRepo) UpsertMulti(ctx context.Context, db database.IDB, s3Storages []*entity.S3Storage,
	conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error {
	if len(s3Storages) == 0 {
		return nil
	}
	query := db.NewInsert().Model(&s3Storages)
	query = bunex.ApplyInsert(query, opts...)
	query = bunex.ApplyUpsert(query, conflictCols, updateCols)

	_, err := query.Exec(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
