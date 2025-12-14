package taskqueue

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/entity"
)

func (q *taskQueue) ScheduleTasks(
	ctx context.Context,
	tasks []*entity.Task,
) (err error) {
	// TODO: add implementation
	return nil
}

func (q *taskQueue) UnscheduleTasks(
	ctx context.Context,
	tasks []*entity.Task,
) error {
	// TODO: add implementation
	return nil
}
