package initializer

import (
	"github.com/localpaas/localpaas/localpaas_app/taskqueue/taskappdeploy"
	"github.com/localpaas/localpaas/localpaas_app/taskqueue/tasktest"
)

type WorkerInitializer struct {
}

func NewWorkerInitializer(
	_ *tasktest.Executor,
	_ *taskappdeploy.Executor,
) *WorkerInitializer {
	return nil
}
