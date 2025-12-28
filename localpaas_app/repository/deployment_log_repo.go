package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
)

type DeploymentLogRepo interface {
	GetByID(ctx context.Context, db database.IDB, id string,
		opts ...bunex.SelectQueryOption) (*entity.DeploymentLog, error)
	List(ctx context.Context, db database.IDB, deploymentID string, paging *basedto.Paging,
		opts ...bunex.SelectQueryOption) ([]*entity.DeploymentLog, *basedto.PagingMeta, error)

	Insert(ctx context.Context, db database.IDB, log *entity.DeploymentLog,
		opts ...bunex.InsertQueryOption) error
	InsertMulti(ctx context.Context, db database.IDB, logs []*entity.DeploymentLog,
		opts ...bunex.InsertQueryOption) error
}

type deploymentLogRepo struct {
}

func NewDeploymentLogRepo() DeploymentLogRepo {
	return &deploymentLogRepo{}
}

func (repo *deploymentLogRepo) GetByID(ctx context.Context, db database.IDB, id string,
	opts ...bunex.SelectQueryOption) (*entity.DeploymentLog, error) {
	log := &entity.DeploymentLog{}
	query := db.NewSelect().Model(log).Where("deployment_log.id = ?", id)
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if log == nil || errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.NewNotFound("DeploymentLog").WithCause(err)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return log, nil
}

func (repo *deploymentLogRepo) List(ctx context.Context, db database.IDB, deploymentID string, paging *basedto.Paging,
	opts ...bunex.SelectQueryOption) ([]*entity.DeploymentLog, *basedto.PagingMeta, error) {
	var logs []*entity.DeploymentLog
	query := db.NewSelect().Model(&logs)
	if deploymentID != "" {
		query = query.Where("deployment_log.deployment_id = ?", deploymentID)
	}
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
	return logs, pagingMeta, nil
}

func (repo *deploymentLogRepo) Insert(ctx context.Context, db database.IDB, log *entity.DeploymentLog,
	opts ...bunex.InsertQueryOption) error {
	return repo.InsertMulti(ctx, db, []*entity.DeploymentLog{log}, opts...)
}

func (repo *deploymentLogRepo) InsertMulti(ctx context.Context, db database.IDB, logs []*entity.DeploymentLog,
	opts ...bunex.InsertQueryOption) error {
	if len(logs) == 0 {
		return nil
	}
	query := db.NewInsert().Model(&logs)
	query = bunex.ApplyInsert(query, opts...)

	_, err := query.Exec(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
