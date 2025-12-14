package taskqueue

import (
	"context"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

func (q *taskQueue) NewTestTaskProcessor() func(taskID string, payload string) error {
	return func(taskID string, payload string) error {
		return q.runTask(context.Background(), taskID, payload, q.doTestTask)
	}
}

func (q *taskQueue) doTestTask(
	ctx context.Context,
	db database.IDB,
	task *entity.Task,
) error {
	// TODO: add implementation
	print(">>>>>>>>>>>>>>>>>>>>> doTestTask ", time.Now().String(), "\n") //nolint
	time.Sleep(10 * time.Second) //nolint
	return nil
}
