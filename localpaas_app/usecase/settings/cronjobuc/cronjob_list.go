package cronjobuc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/entityutil"
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

	input := &cronjobdto.CronJobTransformInput{}
	err = uc.loadReferenceData(ctx, uc.db, resp.Data, input)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := cronjobdto.TransformCronJobs(resp.Data, input)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &cronjobdto.ListCronJobResp{
		Meta: resp.Meta,
		Data: respData,
	}, nil
}

func (uc *CronJobUC) loadReferenceData(
	ctx context.Context,
	db database.IDB,
	cronJobs []*entity.Setting,
	input *cronjobdto.CronJobTransformInput,
) error {
	appIDs := make([]string, 0)
	settingIDs := make([]string, 0, 10) //nolint:mnd
	for _, setting := range cronJobs {
		cronJob := setting.MustAsCronJob()
		if cronJob.App.ID != "" {
			appIDs = append(appIDs, cronJob.App.ID)
		}
		settingIDs = append(settingIDs, cronJob.GetRefSettingIDs()...)
	}

	// Load reference apps
	apps, err := uc.appRepo.ListByIDs(ctx, db, "", gofn.ToSet(appIDs),
		bunex.SelectExcludeColumns(entity.AppDefaultExcludeColumns...),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	input.AppMap = entityutil.SliceToIDMap(apps)

	// Load reference settings
	refSettings, err := uc.settingRepo.ListByIDs(ctx, db, gofn.ToSet(settingIDs), true)
	if err != nil {
		return apperrors.Wrap(err)
	}
	input.RefSettingMap = entityutil.SliceToIDMap(refSettings)

	return nil
}
