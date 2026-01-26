package settings

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/repository"
)

type CreateSettingReq struct {
	BaseSettingReq
	AvailInProjects bool `json:"availableInProjects"`
}

type CreateSettingResp struct {
	Meta *basedto.Meta         `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}

type CreateSettingData struct {
	SettingRepo   repository.SettingRepo
	VerifyingName string
	Version       int

	AfterLoading     func(context.Context, database.Tx, *CreateSettingData) error
	PrepareCreation  func(context.Context, database.Tx, *CreateSettingData, *PersistingSettingCreationData) error
	BeforePersisting func(context.Context, database.Tx, *CreateSettingData, *PersistingSettingCreationData) error
	AfterPersisting  func(context.Context, database.Tx, *CreateSettingData, *PersistingSettingCreationData) error
}

type PersistingSettingCreationData struct {
	Setting *entity.Setting
}

func CreateSetting(
	ctx context.Context,
	db database.IDB,
	req *CreateSettingReq,
	data *CreateSettingData,
) (*CreateSettingResp, error) {
	var persistingData *PersistingSettingCreationData
	err := transaction.Execute(ctx, db, func(db database.Tx) error {
		err := loadSettingForCreation(ctx, db, req, data)
		if err != nil {
			return apperrors.Wrap(err)
		}

		if data.AfterLoading != nil {
			if err := data.AfterLoading(ctx, db, data); err != nil {
				return apperrors.Wrap(err)
			}
		}

		persistingData = &PersistingSettingCreationData{}
		prepareSettingCreation(req, data, persistingData)

		if data.PrepareCreation != nil {
			if err := data.PrepareCreation(ctx, db, data, persistingData); err != nil {
				return apperrors.Wrap(err)
			}
		}

		if data.BeforePersisting != nil {
			if err := data.BeforePersisting(ctx, db, data, persistingData); err != nil {
				return apperrors.Wrap(err)
			}
		}

		err = persistSettingCreation(ctx, db, data, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		if data.AfterPersisting != nil {
			if err := data.AfterPersisting(ctx, db, data, persistingData); err != nil {
				return apperrors.Wrap(err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &CreateSettingResp{
		Data: &basedto.ObjectIDResp{ID: persistingData.Setting.ID},
	}, nil
}

func loadSettingForCreation(
	ctx context.Context,
	db database.IDB,
	req *CreateSettingReq,
	data *CreateSettingData,
) error {
	// Verify that the name is available to use
	err := checkNameConflict(ctx, db, data.SettingRepo, &req.BaseSettingReq, data.VerifyingName)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func prepareSettingCreation(
	req *CreateSettingReq,
	data *CreateSettingData,
	persistingData *PersistingSettingCreationData,
) {
	timeNow := timeutil.NowUTC()
	setting := &entity.Setting{
		ID:              gofn.Must(ulid.NewStringULID()),
		Type:            req.Type,
		Status:          base.SettingStatusActive,
		Name:            data.VerifyingName,
		ObjectID:        req.ObjectID,
		AvailInProjects: gofn.If(req.Scope != base.SettingScopeGlobal, false, req.AvailInProjects),
		Version:         data.Version,
		CreatedAt:       timeNow,
		UpdatedAt:       timeNow,
	}
	persistingData.Setting = setting
}

func persistSettingCreation(
	ctx context.Context,
	db database.IDB,
	data *CreateSettingData,
	persistingData *PersistingSettingCreationData,
) error {
	err := data.SettingRepo.Upsert(ctx, db, persistingData.Setting,
		entity.SettingUpsertingConflictCols, entity.SettingUpsertingUpdateCols)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
