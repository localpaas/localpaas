package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
)

type LockRepo interface {
	GetByID(ctx context.Context, db database.IDB, id string,
		opts ...bunex.SelectQueryOption) (*entity.Lock, error)

	Upsert(ctx context.Context, db database.IDB, lock *entity.Lock, conflictCols, updateCols []string,
		opts ...bunex.InsertQueryOption) error
	UpsertMulti(ctx context.Context, db database.IDB, locks []*entity.Lock, conflictCols, updateCols []string,
		opts ...bunex.InsertQueryOption) error
}

type lockRepo struct {
}

func NewLockRepo() LockRepo {
	return &lockRepo{}
}

func (repo *lockRepo) GetByID(ctx context.Context, db database.IDB, id string,
	opts ...bunex.SelectQueryOption) (*entity.Lock, error) {
	lock := &entity.Lock{}
	query := db.NewSelect().Model(lock).Where("lock.id = ?", id)
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if lock == nil || errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.NewNotFound("Lock").WithCause(err)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return lock, nil
}

func (repo *lockRepo) Upsert(ctx context.Context, db database.IDB, lock *entity.Lock,
	conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error {
	return repo.UpsertMulti(ctx, db, []*entity.Lock{lock}, conflictCols, updateCols, opts...)
}

func (repo *lockRepo) UpsertMulti(ctx context.Context, db database.IDB, locks []*entity.Lock,
	conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error {
	if len(locks) == 0 {
		return nil
	}
	query := db.NewInsert().Model(&locks)
	query = bunex.ApplyInsert(query, opts...)
	query = bunex.ApplyUpsert(query, conflictCols, updateCols)

	_, err := query.Exec(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
