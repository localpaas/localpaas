package cronjobuc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/cronjobuc/cronjobdto"
)

func (uc *CronJobUC) UpdateCronJob(
	ctx context.Context,
	auth *basedto.Auth,
	req *cronjobdto.UpdateCronJobReq,
) (*cronjobdto.UpdateCronJobResp, error) {
	req.Type = currentSettingType
	unscheduleCurrentTasks := false
	_, err := settings.UpdateSetting(ctx, uc.db, &req.UpdateSettingReq, &settings.UpdateSettingData{
		SettingRepo:   uc.settingRepo,
		VerifyingName: req.Name,
		AfterLoading: func(ctx context.Context, db database.Tx, data *settings.UpdateSettingData) error {
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
			data *settings.UpdateSettingData,
			pData *settings.PersistingSettingData,
		) error {
			pData.Setting.Name = gofn.Coalesce(req.Name, pData.Setting.Name)
			err := pData.Setting.SetData(req.ToEntity())
			if err != nil {
				return apperrors.Wrap(err)
			}
			return nil
		},
		AfterPersisting: func(ctx context.Context, db database.Tx, data *settings.UpdateSettingData,
			pData *settings.PersistingSettingData) error {
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
