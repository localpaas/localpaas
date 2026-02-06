package initializer

import (
	"github.com/localpaas/localpaas/localpaas_app/taskqueue/taskappdeploy"
	"github.com/localpaas/localpaas/localpaas_app/taskqueue/taskappnotification"
	"github.com/localpaas/localpaas/localpaas_app/taskqueue/taskcronjobexec"
	"github.com/localpaas/localpaas/localpaas_app/taskqueue/tasktest"
)

type WorkerInitializer struct {
}

// NOTE: these injections are required to make the task executor be available
func NewWorkerInitializer(
	_ *tasktest.Executor,
	_ *taskappdeploy.Executor,
	_ *taskappnotification.Executor,
	_ *taskcronjobexec.Executor,
) *WorkerInitializer {
	return nil
}
