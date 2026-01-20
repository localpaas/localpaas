package cronjobuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/cronjobuc/cronjobdto"
)

const (
	currentSettingType    = base.SettingTypeCronJob
	currentSettingVersion = entity.CurrentCronJobVersion
)

func (uc *CronJobUC) CreateCronJob(
	ctx context.Context,
	auth *basedto.Auth,
	req *cronjobdto.CreateCronJobReq,
) (*cronjobdto.CreateCronJobResp, error) {
	req.Type = currentSettingType
	resp, err := providers.CreateSetting(ctx, uc.db, &req.CreateSettingReq, &providers.CreateSettingData{
		SettingRepo:   uc.settingRepo,
		VerifyingName: req.Name,
		Version:       currentSettingVersion,
		PrepareCreation: func(
			ctx context.Context,
			db database.Tx,
			data *providers.CreateSettingData,
			pData *providers.PersistingSettingCreationData,
		) error {
			pData.Setting.Kind = string(req.Kind)
			cronJob := &entity.CronJob{
				Cron:        req.Cron,
				InitialTime: timeutil.NowUTC(),
				Priority:    req.Priority,
				MaxRetry:    req.MaxRetry,
				RetryDelay:  req.RetryDelay,
				Timeout:     req.Timeout,
				Command:     req.Command,
			}
			// Parse the cron expression to make sure it's valid
			_, err := cronJob.ParseCron()
			if err != nil {
				return apperrors.Wrap(err)
			}
			if err = pData.Setting.SetData(cronJob); err != nil {
				return apperrors.Wrap(err)
			}
			return nil
		},
		AfterPersisting: func(ctx context.Context, db database.Tx, data *providers.CreateSettingData,
			pData *providers.PersistingSettingCreationData) error {
			err := uc.taskQueue.ScheduleTasksForCronJob(ctx, db, pData.Setting, false)
			if err != nil {
				return apperrors.Wrap(err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &cronjobdto.CreateCronJobResp{
		Data: resp.Data,
	}, nil
}
