package cronjobuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/cronjobuc/cronjobdto"
)

func (uc *CronJobUC) DeleteCronJob(
	ctx context.Context,
	auth *basedto.Auth,
	req *cronjobdto.DeleteCronJobReq,
) (*cronjobdto.DeleteCronJobResp, error) {
	var jobData *deleteCronJobData
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		jobData = &deleteCronJobData{}
		err := uc.loadCronJobDataForDelete(ctx, db, req, jobData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingCronJobData{}
		uc.prepareDeletingCronJob(jobData, persistingData)

		err = uc.persistData(ctx, db, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		err = uc.taskQueue.ScheduleTasksForCronJob(ctx, db, jobData.Setting, jobData.UnscheduleCurrentTasks)
		if err != nil {
			return apperrors.Wrap(err)
		}

		return nil
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &cronjobdto.DeleteCronJobResp{}, nil
}

type deleteCronJobData struct {
	Setting                *entity.Setting
	UnscheduleCurrentTasks bool
}

func (uc *CronJobUC) loadCronJobDataForDelete(
	ctx context.Context,
	db database.IDB,
	req *cronjobdto.DeleteCronJobReq,
	data *deleteCronJobData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeCronJob, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	data.UnscheduleCurrentTasks = true
	return nil
}

func (uc *CronJobUC) prepareDeletingCronJob(
	data *deleteCronJobData,
	persistingData *persistingCronJobData,
) {
	setting := data.Setting
	setting.DeletedAt = timeutil.NowUTC()
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
