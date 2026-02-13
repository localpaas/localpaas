package healthcheckuc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/entityutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/healthcheckuc/healthcheckdto"
)

func (uc *HealthcheckUC) ListHealthcheck(
	ctx context.Context,
	auth *basedto.Auth,
	req *healthcheckdto.ListHealthcheckReq,
) (*healthcheckdto.ListHealthcheckResp, error) {
	req.Type = currentSettingType
	resp, err := uc.ListSetting(ctx, auth, &req.ListSettingReq, &settings.ListSettingData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	input := &healthcheckdto.HealthcheckTransformInput{}
	err = uc.loadReferenceData(ctx, uc.DB, resp.Data, input)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := healthcheckdto.TransformHealthchecks(resp.Data, input)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &healthcheckdto.ListHealthcheckResp{
		Meta: resp.Meta,
		Data: respData,
	}, nil
}

func (uc *HealthcheckUC) loadReferenceData(
	ctx context.Context,
	db database.IDB,
	cronJobs []*entity.Setting,
	input *healthcheckdto.HealthcheckTransformInput,
) error {
	settingIDs := make([]string, 0, 10) //nolint:mnd
	for _, setting := range cronJobs {
		cronJob := setting.MustAsHealthcheck()
		settingIDs = append(settingIDs, cronJob.GetRefSettingIDs()...)
	}

	// Load reference settings
	refSettings, err := uc.SettingRepo.ListByIDs(ctx, db, gofn.ToSet(settingIDs), true)
	if err != nil {
		return apperrors.Wrap(err)
	}
	input.RefSettingMap = entityutil.SliceToIDMap(refSettings)

	return nil
}
