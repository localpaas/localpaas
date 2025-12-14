package repository

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
)

type UpdatingTaskRepo interface {
	Upsert(ctx context.Context, db database.IDB, task *entity.UpdatingTask,
		conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error

	Delete(ctx context.Context, db database.IDB, task *entity.UpdatingTask,
		opts ...bunex.DeleteQueryOption) error
}

type updatingTaskRepo struct {
}

func NewUpdatingTaskRepo() UpdatingTaskRepo {
	return &updatingTaskRepo{}
}

func (repo *updatingTaskRepo) Upsert(ctx context.Context, db database.IDB, task *entity.UpdatingTask,
	conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error {
	query := db.NewInsert().Model(task)
	query = bunex.ApplyInsert(query, opts...)
	query = bunex.ApplyUpsert(query, conflictCols, updateCols)

	_, err := query.Exec(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (repo *updatingTaskRepo) Delete(ctx context.Context, db database.IDB, task *entity.UpdatingTask,
	opts ...bunex.DeleteQueryOption) error {
	query := db.NewDelete().Model(task).WherePK()
	query = bunex.ApplyDelete(query, opts...)

	_, err := query.Exec(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
