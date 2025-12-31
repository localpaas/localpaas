package cronjobuc

import (
	"context"
	"strings"

	"github.com/tiendc/gofn"

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

func (uc *CronJobUC) UpdateCronJob(
	ctx context.Context,
	auth *basedto.Auth,
	req *cronjobdto.UpdateCronJobReq,
) (*cronjobdto.UpdateCronJobResp, error) {
	var jobData *updateCronJobData
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		jobData = &updateCronJobData{}
		err := uc.loadCronJobDataForUpdate(ctx, db, req, jobData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingCronJobData{}
		uc.prepareUpdatingCronJob(req.CronJobBaseReq, jobData, persistingData)

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

	return &cronjobdto.UpdateCronJobResp{}, nil
}

type updateCronJobData struct {
	Setting                *entity.Setting
	UnscheduleCurrentTasks bool
}

func (uc *CronJobUC) loadCronJobDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *cronjobdto.UpdateCronJobReq,
	data *updateCronJobData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeCronJob, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if req.UpdateVer != setting.UpdateVer {
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}
	data.Setting = setting

	// If name changes, validate the new one
	name := gofn.Coalesce(req.Name, setting.Name)
	if name != "" && !strings.EqualFold(setting.Name, name) {
		conflictSetting, _ := uc.settingRepo.GetByName(ctx, db, base.SettingTypeCronJob, name, false)
		if conflictSetting != nil {
			return apperrors.NewAlreadyExist("CronJob").
				WithMsgLog("cron job '%s' already exists", conflictSetting.Name)
		}
	}

	job, err := setting.AsCronJob()
	if err != nil {
		return apperrors.Wrap(err)
	}

	if req.Cron != job.Cron {
		data.UnscheduleCurrentTasks = true
	}

	return nil
}

func (uc *CronJobUC) prepareUpdatingCronJob(
	req *cronjobdto.CronJobBaseReq,
	data *updateCronJobData,
	persistingData *persistingCronJobData,
) {
	timeNow := timeutil.NowUTC()
	setting := data.Setting
	setting.UpdateVer++
	setting.UpdatedAt = timeNow
	setting.Name = gofn.Coalesce(req.Name, setting.Name)

	cronJob := &entity.CronJob{
		Cron:         req.Cron,
		InitialTime:  timeNow,
		Priority:     req.Priority,
		MaxRetry:     req.MaxRetry,
		RetryDelayMs: req.RetryDelayMs,
		TimeoutMs:    req.TimeoutMs,
		Command:      req.Command,
	}
	setting.MustSetData(cronJob)
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
