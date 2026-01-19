package providers

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/strutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/repository"
)

type CreateSettingReq struct {
	Type       base.SettingType `json:"-"`
	ProjectID  string           `json:"-"`
	AppID      string           `json:"-"`
	GlobalOnly bool             `json:"-"`
}

type CreateSettingResp struct {
	Meta *basedto.BaseMeta     `json:"meta"`
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
) (err error) {
	// Verify that the name is available to use
	if data.VerifyingName != "" {
		conflictSetting, _ := data.SettingRepo.GetByNameEx(ctx, db, req.Type,
			req.ProjectID, req.AppID, data.VerifyingName, false,
			bunex.SelectWhereIf(req.GlobalOnly, "setting.object_id IS NULL"),
		)
		if conflictSetting != nil {
			return apperrors.NewAlreadyExist(strutil.ToPascalCase(string(req.Type))).
				WithMsgLog("%s '%s' already exists", req.Type, conflictSetting.Name)
		}
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
		ID:        gofn.Must(ulid.NewStringULID()),
		Type:      req.Type,
		Status:    base.SettingStatusActive,
		Name:      data.VerifyingName,
		ObjectID:  gofn.Coalesce(req.ProjectID, req.AppID),
		Version:   data.Version,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
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
