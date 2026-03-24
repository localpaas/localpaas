package taskcronjobexec

import (
	"context"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
)

type sysBackupTaskData struct {
	*taskData
	TaskOutput     *entity.TaskSystemBackupOutput
	BackupDir      string
	BackupFileDir  string
	BackupFileName string
	TimeNow        time.Time
}

func (e *Executor) cronExecSystemBackup(
	ctx context.Context,
	db database.IDB,
	data *taskData,
) error {
	setting := data.RefObjects.RefSettings[data.CronJob.TargetSetting.ID]
	if setting == nil {
		return apperrors.NewNotFound("System backup settings")
	}
	sysBackup := setting.MustAsSystemBackup()

	taskData := &sysBackupTaskData{
		taskData: data,
		TaskOutput: &entity.TaskSystemBackupOutput{
			DBBackup: &entity.DBBackupOutput{},
		},
		TimeNow: timeutil.NowUTC(),
	}

	// Backup DB
	err := e.sysDBBackup(ctx, db, sysBackup, taskData)

	// Assign back the result output
	data.Task.MustSetOutput(taskData.TaskOutput)

	return err
}
