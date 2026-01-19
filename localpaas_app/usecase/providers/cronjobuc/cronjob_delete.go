package cronjobuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/cronjobuc/cronjobdto"
)

func (uc *CronJobUC) DeleteCronJob(
	ctx context.Context,
	auth *basedto.Auth,
	req *cronjobdto.DeleteCronJobReq,
) (*cronjobdto.DeleteCronJobResp, error) {
	req.Type = currentSettingType
	_, err := providers.DeleteSetting(ctx, uc.db, &req.DeleteSettingReq, &providers.DeleteSettingData{
		SettingRepo:              uc.settingRepo,
		ProjectSharedSettingRepo: uc.projectSharedSettingRepo,
		AfterPersisting: func(ctx context.Context, db database.Tx, data *providers.DeleteSettingData,
			_ *providers.PersistingSettingDeletionData) error {
			err := uc.taskQueue.ScheduleTasksForCronJob(ctx, db, data.Setting, true)
			if err != nil {
				return apperrors.Wrap(err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &cronjobdto.DeleteCronJobResp{}, nil
}
