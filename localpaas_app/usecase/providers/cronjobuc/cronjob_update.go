package cronjobuc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/cronjobuc/cronjobdto"
)

func (uc *CronJobUC) UpdateCronJob(
	ctx context.Context,
	auth *basedto.Auth,
	req *cronjobdto.UpdateCronJobReq,
) (*cronjobdto.UpdateCronJobResp, error) {
	req.Type = currentSettingType
	unscheduleCurrentTasks := false
	_, err := providers.UpdateSetting(ctx, uc.db, &req.UpdateSettingReq, &providers.UpdateSettingData{
		SettingRepo:   uc.settingRepo,
		VerifyingName: req.Name,
		AfterLoading: func(ctx context.Context, db database.Tx, data *providers.UpdateSettingData) error {
			job, err := data.Setting.AsCronJob()
			if err != nil {
				return apperrors.Wrap(err)
			}
			unscheduleCurrentTasks = req.Cron != job.Cron
			return nil
		},
		PrepareUpdate: func(
			ctx context.Context,
			db database.Tx,
			data *providers.UpdateSettingData,
			pData *providers.PersistingSettingData,
		) error {
			pData.Setting.Name = gofn.Coalesce(req.Name, pData.Setting.Name)
			err := pData.Setting.SetData(&entity.CronJob{
				Cron:        req.Cron,
				InitialTime: timeutil.NowUTC(),
				Priority:    req.Priority,
				MaxRetry:    req.MaxRetry,
				RetryDelay:  req.RetryDelay,
				Timeout:     req.Timeout,
				Command:     req.Command,
			})
			if err != nil {
				return apperrors.Wrap(err)
			}
			return nil
		},
		AfterPersisting: func(ctx context.Context, db database.Tx, data *providers.UpdateSettingData,
			pData *providers.PersistingSettingData) error {
			err := uc.taskQueue.ScheduleTasksForCronJob(ctx, db, data.Setting, unscheduleCurrentTasks)
			if err != nil {
				return apperrors.Wrap(err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &cronjobdto.UpdateCronJobResp{}, nil
}
