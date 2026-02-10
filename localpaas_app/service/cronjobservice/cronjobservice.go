package cronjobservice

import (
	"context"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/envvarservice"
)

type CronJobService interface {
	BuildCommandEnv(ctx context.Context, db database.IDB, app *entity.App, cronJob *entity.CronJob) (
		res []*envvarservice.EnvVar, err error)

	CreateCronJobTask(job *entity.Setting, runAt, timeNow time.Time) (*entity.Task, error)
}

func NewCronJobService(
	settingRepo repository.SettingRepo,
	envVarService envvarservice.EnvVarService,
) CronJobService {
	return &cronJobService{
		settingRepo:   settingRepo,
		envVarService: envVarService,
	}
}

type cronJobService struct {
	settingRepo   repository.SettingRepo
	envVarService envvarservice.EnvVarService
}
