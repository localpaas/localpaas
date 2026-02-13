package cronjobuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/cronjobuc/cronjobdto"
)

func (uc *CronJobUC) GetCronJob(
	ctx context.Context,
	auth *basedto.Auth,
	req *cronjobdto.GetCronJobReq,
) (*cronjobdto.GetCronJobResp, error) {
	req.Type = currentSettingType
	setting, err := uc.GetSetting(ctx, auth, &req.GetSettingReq, &settings.GetSettingData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	input := &cronjobdto.CronJobTransformInput{}
	err = uc.loadReferenceData(ctx, uc.DB, []*entity.Setting{setting}, input)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := cronjobdto.TransformCronJob(setting, input)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &cronjobdto.GetCronJobResp{
		Data: resp,
	}, nil
}
