package taskservice

import (
	"context"
	"errors"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
)

const (
	maxTryLock = 3
)

func (s *taskService) CreateDBLock(
	ctx context.Context,
	db database.Tx,
	id string,
	selectFor string,
) (*entity.Lock, error) {
	return s.tryCreateDBLock(ctx, db, id, selectFor, 1)
}

func (s *taskService) tryCreateDBLock(
	ctx context.Context,
	db database.Tx,
	id string,
	selectFor string,
	try int,
) (*entity.Lock, error) {
	if try >= maxTryLock {
		return nil, apperrors.Wrap(apperrors.ErrActionFailed)
	}
	if selectFor == "" {
		selectFor = "UPDATE"
	}
	lock, err := s.lockRepo.GetByID(ctx, db, id,
		bunex.SelectFor(selectFor),
	)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, apperrors.Wrap(err)
	}

	if lock != nil {
		return lock, nil
	}

	// Insert a new row for the lock using the default DB (not the transaction one)
	lock = &entity.Lock{ID: id}
	err = s.lockRepo.Upsert(ctx, s.db, lock, entity.LockUpsertingConflictCols, entity.LockUpsertingUpdateCols)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	try++
	return s.tryCreateDBLock(ctx, db, id, selectFor, try)
}
