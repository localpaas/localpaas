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
	resp, err := uc.GetSetting(ctx, auth, &req.GetSettingReq, &settings.GetSettingData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	input := &cronjobdto.CronJobTransformInput{
		RefObjects: resp.RefObjects,
	}
	err = uc.loadReferenceData(ctx, uc.DB, []*entity.Setting{resp.Data}, input)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := cronjobdto.TransformCronJob(resp.Data, input)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &cronjobdto.GetCronJobResp{
		Data: respData,
	}, nil
}
