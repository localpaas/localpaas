package taskcronjobexec

import (
	"context"
	"errors"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type sysCleanupTaskData struct {
	*taskData
	TaskOutput *entity.TaskSystemCleanupOutput
}

func (e *Executor) cronExecSystemCleanup(
	ctx context.Context,
	db database.IDB,
	data *taskData,
) error {
	cronJob := data.CronJob.MustAsCronJob()
	setting := data.RefObjects.RefSettings[cronJob.TargetSetting.ID]
	if setting == nil {
		return apperrors.NewNotFound("System cleanup settings")
	}
	sysCleanup := setting.MustAsSystemCleanup()

	taskData := &sysCleanupTaskData{
		taskData: data,
		TaskOutput: &entity.TaskSystemCleanupOutput{
			DBCleanup:      &entity.DBCleanupOutput{},
			ClusterCleanup: &entity.ClusterCleanupOutput{},
			FileCleanup:    &entity.FileCleanupOutput{},
		},
	}

	// Cleanup DB objects
	err1 := e.sysCleanupDB(ctx, db, sysCleanup.DBObjectRetention, taskData)

	// Cleanup unused cluster data (docker)
	err2 := e.sysCleanupCluster(ctx, sysCleanup.ClusterCleanup, taskData)

	// Cleanup orphaned files
	err3 := e.sysCleanupFiles(ctx, taskData)

	// Assign back the result output
	data.Task.MustSetOutput(taskData.TaskOutput)

	return errors.Join(err1, err2, err3)
}
