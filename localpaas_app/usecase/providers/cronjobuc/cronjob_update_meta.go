package cronjobuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/cronjobuc/cronjobdto"
)

func (uc *CronJobUC) UpdateCronJobMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *cronjobdto.UpdateCronJobMetaReq,
) (*cronjobdto.UpdateCronJobMetaResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		jobData := &updateCronJobData{}
		err := uc.loadCronJobDataForUpdateMeta(ctx, db, req, jobData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		uc.prepareUpdatingCronJobMeta(req, jobData)
		err = uc.persistCronJobMeta(ctx, db, jobData)
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

	return &cronjobdto.UpdateCronJobMetaResp{}, nil
}

func (uc *CronJobUC) loadCronJobDataForUpdateMeta(
	ctx context.Context,
	db database.IDB,
	req *cronjobdto.UpdateCronJobMetaReq,
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

	data.UnscheduleCurrentTasks = true

	return nil
}

func (uc *CronJobUC) prepareUpdatingCronJobMeta(
	req *cronjobdto.UpdateCronJobMetaReq,
	data *updateCronJobData,
) {
	timeNow := timeutil.NowUTC()
	setting := data.Setting
	setting.UpdateVer++
	setting.UpdatedAt = timeNow

	if req.Status != nil {
		setting.Status = *req.Status
	}
	if req.ExpireAt != nil {
		setting.ExpireAt = *req.ExpireAt
	}
}

func (uc *CronJobUC) persistCronJobMeta(
	ctx context.Context,
	db database.IDB,
	data *updateCronJobData,
) error {
	err := uc.settingRepo.Update(ctx, db, data.Setting,
		bunex.UpdateColumns("update_ver", "updated_at", "status", "expire_at"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
