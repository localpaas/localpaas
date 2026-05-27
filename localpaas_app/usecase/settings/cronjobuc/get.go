package cronjobuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/cronjobuc/cronjobdto"
)

const (
	cronJobNextRunsCalculation = 5
)

func (uc *UC) GetCronJob(
	ctx context.Context,
	auth *basedto.Auth,
	req *cronjobdto.GetCronJobReq,
) (*cronjobdto.GetCronJobResp, error) {
	req.Type = currentSettingType
	resp, err := uc.GetSetting(ctx, auth, &req.GetSettingReq, &settings.GetSettingData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := cronjobdto.TransformCronJob(resp.Data, resp.RefObjects)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	// Return few next runs of the job
	sched := resp.Data.MustAsCronJob().Schedule
	respData.NextRuns, _ = sched.CalcNextRuns(timeutil.NowUTC(), cronJobNextRunsCalculation)

	return &cronjobdto.GetCronJobResp{
		Data: respData,
	}, nil
}
