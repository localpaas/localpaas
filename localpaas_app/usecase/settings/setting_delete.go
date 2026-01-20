package settings

import (
	"context"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/repository"
)

type DeleteSettingReq struct {
	ID         string           `json:"-"`
	Type       base.SettingType `json:"-"`
	ProjectID  string           `json:"-"`
	AppID      string           `json:"-"`
	GlobalOnly bool             `json:"-"`
}

func (req *DeleteSettingReq) Validate() (validators []vld.Validator) {
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return
}

type DeleteSettingResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}

type DeleteSettingData struct {
	Setting              *entity.Setting
	ProjectSharedSetting *entity.ProjectSharedSetting

	SettingRepo              repository.SettingRepo
	ProjectSharedSettingRepo repository.ProjectSharedSettingRepo
	ExtraLoadOpts            []bunex.SelectQueryOption

	AfterLoading     func(context.Context, database.Tx, *DeleteSettingData) error
	BeforePersisting func(context.Context, database.Tx, *DeleteSettingData, *PersistingSettingDeletionData) error
	AfterPersisting  func(context.Context, database.Tx, *DeleteSettingData, *PersistingSettingDeletionData) error
}

type PersistingSettingDeletionData struct {
	Setting              *entity.Setting
	ProjectSharedSetting *entity.ProjectSharedSetting
}

func DeleteSetting(
	ctx context.Context,
	db database.IDB,
	req *DeleteSettingReq,
	data *DeleteSettingData,
) (*DeleteSettingResp, error) {
	err := transaction.Execute(ctx, db, func(db database.Tx) error {
		err := loadSettingForDeletion(ctx, db, req, data)
		if err != nil {
			return apperrors.Wrap(err)
		}

		if data.AfterLoading != nil {
			if err := data.AfterLoading(ctx, db, data); err != nil {
				return apperrors.Wrap(err)
			}
		}

		persistingData := &PersistingSettingDeletionData{}
		prepareSettingDeletion(req, data, persistingData)
		if data.BeforePersisting != nil {
			if err := data.BeforePersisting(ctx, db, data, persistingData); err != nil {
				return apperrors.Wrap(err)
			}
		}

		err = persistSettingDeletion(ctx, db, data, persistingData)
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

	return &DeleteSettingResp{}, nil
}

func loadSettingForDeletion(
	ctx context.Context,
	db database.IDB,
	req *DeleteSettingReq,
	data *DeleteSettingData,
) (err error) {
	loadOpts := []bunex.SelectQueryOption{
		bunex.SelectFor("UPDATE OF setting"),
		bunex.SelectWhereIf(req.GlobalOnly, "setting.object_id IS NULL"),
	}
	loadOpts = append(loadOpts, data.ExtraLoadOpts...)

	setting, err := data.SettingRepo.GetByIDEx(ctx, db, req.Type, req.ProjectID, req.AppID, req.ID,
		false, loadOpts...)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	if setting.ObjectID == "" && req.ProjectID != "" {
		data.ProjectSharedSetting, err = data.ProjectSharedSettingRepo.Get(ctx, db, req.ProjectID, req.ID)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	return nil
}

func prepareSettingDeletion(
	_ *DeleteSettingReq,
	data *DeleteSettingData,
	persistingData *PersistingSettingDeletionData,
) {
	timeNow := timeutil.NowUTC()
	if data.ProjectSharedSetting != nil {
		data.ProjectSharedSetting.DeletedAt = timeNow
		persistingData.ProjectSharedSetting = data.ProjectSharedSetting
	} else {
		data.Setting.UpdateVer++
		data.Setting.DeletedAt = timeNow
		persistingData.Setting = data.Setting
	}
}

func persistSettingDeletion(
	ctx context.Context,
	db database.IDB,
	data *DeleteSettingData,
	persistingData *PersistingSettingDeletionData,
) (err error) {
	if data.ProjectSharedSetting != nil {
		err = data.ProjectSharedSettingRepo.Update(ctx, db, persistingData.ProjectSharedSetting,
			bunex.UpdateColumns("deleted_at"),
		)
	} else {
		err = data.SettingRepo.Update(ctx, db, persistingData.Setting,
			bunex.UpdateColumns("deleted_at"),
		)
	}
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
