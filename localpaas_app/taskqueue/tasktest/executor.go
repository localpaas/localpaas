package tasktest

import (
	"context"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/taskqueue"
)

type Executor struct {
	logger      logging.Logger
	settingRepo repository.SettingRepo
}

func NewExecutor(
	taskQueue taskqueue.TaskQueue,
	logger logging.Logger,
	settingRepo repository.SettingRepo,
) *Executor {
	p := &Executor{
		logger:      logger,
		settingRepo: settingRepo,
	}
	taskQueue.RegisterExecutor(base.TaskTypeTest, p.execute)
	return p
}

func (p *Executor) execute(
	ctx context.Context,
	db database.Tx,
	task *entity.Task,
) error {
	// TODO: add implementation
	print(">>>>>>>>>>>>>>>>>>>>> doTestTask ", time.Now().String(), task.Job.Name, "\n") //nolint
	time.Sleep(3 * time.Second)                                                          //nolint
	return nil
}
