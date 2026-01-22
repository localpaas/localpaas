package cronjobuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/cronjobuc/cronjobdto"
)

func (uc *CronJobUC) ListCronJob(
	ctx context.Context,
	auth *basedto.Auth,
	req *cronjobdto.ListCronJobReq,
) (*cronjobdto.ListCronJobResp, error) {
	req.Type = currentSettingType
	resp, err := settings.ListSetting(ctx, uc.db, auth, &req.ListSettingReq, &settings.ListSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := cronjobdto.TransformCronJobs(resp.Data, req.ObjectID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &cronjobdto.ListCronJobResp{
		Meta: resp.Meta,
		Data: respData,
	}, nil
}
