package taskqueue

import (
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/infra/gocronqueue"
)

func (q *taskQueue) getTaskMap() map[base.TaskType]gocronqueue.TaskProcessorFunc {
	return map[base.TaskType]gocronqueue.TaskProcessorFunc{
		base.TaskTypeTest: q.NewTaskTestProcessor(),
	}
}
