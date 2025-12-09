package cronjobuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/cronjobuc/cronjobdto"
)

func (uc *CronJobUC) GetCronJob(
	ctx context.Context,
	auth *basedto.Auth,
	req *cronjobdto.GetCronJobReq,
) (*cronjobdto.GetCronJobResp, error) {
	setting, err := uc.settingRepo.GetByID(ctx, uc.db, base.SettingTypeCronJob, req.ID, false)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := cronjobdto.TransformCronJob(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &cronjobdto.GetCronJobResp{
		Data: resp,
	}, nil
}
