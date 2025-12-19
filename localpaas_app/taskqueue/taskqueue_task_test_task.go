package taskqueue

import (
	"context"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

func (q *taskQueue) NewTaskTestProcessor() func(taskID string, payload string) (time.Time, error) {
	return func(taskID string, payload string) (time.Time, error) {
		return q.runTask(context.Background(), taskID, payload, q.doTestTask)
	}
}

func (q *taskQueue) doTestTask(
	ctx context.Context,
	db database.Tx,
	task *entity.Task,
) error {
	// TODO: add implementation
	print(">>>>>>>>>>>>>>>>>>>>> doTestTask ", time.Now().String(), task.Job.Name, "\n") //nolint
	time.Sleep(3 * time.Second)                                                          //nolint
	return nil
}
