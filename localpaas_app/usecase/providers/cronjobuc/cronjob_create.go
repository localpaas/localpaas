package cronjobuc

import (
	"context"
	"errors"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/cronjobuc/cronjobdto"
)

func (uc *CronJobUC) CreateCronJob(
	ctx context.Context,
	auth *basedto.Auth,
	req *cronjobdto.CreateCronJobReq,
) (*cronjobdto.CreateCronJobResp, error) {
	jobData := &createCronJobData{}
	err := uc.loadCronJobData(ctx, uc.db, req, jobData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	persistingData := &persistingCronJobData{}
	uc.preparePersistingCronJob(req, jobData, persistingData)

	err = transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	createdItem := persistingData.UpsertingSettings[0]
	return &cronjobdto.CreateCronJobResp{
		Data: &basedto.ObjectIDResp{ID: createdItem.ID},
	}, nil
}

type createCronJobData struct {
}

func (uc *CronJobUC) loadCronJobData(
	ctx context.Context,
	db database.IDB,
	req *cronjobdto.CreateCronJobReq,
	_ *createCronJobData,
) error {
	setting, err := uc.settingRepo.GetByName(ctx, db, base.SettingTypeCronJob, req.Name, false)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if setting != nil {
		return apperrors.NewAlreadyExist("CronJob").
			WithMsgLog("cron job '%s' already exists", req.Name)
	}

	return nil
}

type persistingCronJobData struct {
	settingservice.PersistingSettingData
}

func (uc *CronJobUC) preparePersistingCronJob(
	req *cronjobdto.CreateCronJobReq,
	_ *createCronJobData,
	persistingData *persistingCronJobData,
) {
	timeNow := timeutil.NowUTC()
	setting := &entity.Setting{
		ID:        gofn.Must(ulid.NewStringULID()),
		Type:      base.SettingTypeCronJob,
		Status:    base.SettingStatusActive,
		Kind:      string(base.SettingTypeCronJob),
		Name:      req.Name,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	cronJob := &entity.CronJob{
		Cron:        req.Cron,
		InitialTime: timeNow,
		Command:     req.Command,
	}
	setting.MustSetData(cronJob)

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}

func (uc *CronJobUC) persistData(
	ctx context.Context,
	db database.IDB,
	persistingData *persistingCronJobData,
) error {
	err := uc.settingService.PersistSettingData(ctx, db, &persistingData.PersistingSettingData)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
