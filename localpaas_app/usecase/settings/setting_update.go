package settings

import (
	"context"
	"strings"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/strutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/repository"
)

type UpdateSettingReq struct {
	ID         string           `json:"-"`
	Type       base.SettingType `json:"-"`
	ProjectID  string           `json:"-"`
	AppID      string           `json:"-"`
	GlobalOnly bool             `json:"-"`
	UpdateVer  int              `json:"updateVer"`
}

type UpdateSettingResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}

type UpdateSettingData struct {
	Setting *entity.Setting

	SettingRepo   repository.SettingRepo
	VerifyingName string
	ExtraLoadOpts []bunex.SelectQueryOption

	AfterLoading     func(context.Context, database.Tx, *UpdateSettingData) error
	PrepareUpdate    func(context.Context, database.Tx, *UpdateSettingData, *PersistingSettingData) error
	BeforePersisting func(context.Context, database.Tx, *UpdateSettingData, *PersistingSettingData) error
	AfterPersisting  func(context.Context, database.Tx, *UpdateSettingData, *PersistingSettingData) error
}

type PersistingSettingData struct {
	Setting *entity.Setting
}

func UpdateSetting(
	ctx context.Context,
	db database.IDB,
	req *UpdateSettingReq,
	data *UpdateSettingData,
) (*UpdateSettingResp, error) {
	err := transaction.Execute(ctx, db, func(db database.Tx) error {
		err := loadSettingForUpdate(ctx, db, req, data)
		if err != nil {
			return apperrors.Wrap(err)
		}

		if data.AfterLoading != nil {
			if err := data.AfterLoading(ctx, db, data); err != nil {
				return apperrors.Wrap(err)
			}
		}

		persistingData := &PersistingSettingData{}
		prepareSettingUpdate(req, data, persistingData)

		if data.PrepareUpdate != nil {
			if err := data.PrepareUpdate(ctx, db, data, persistingData); err != nil {
				return apperrors.Wrap(err)
			}
		}

		if data.BeforePersisting != nil {
			if err := data.BeforePersisting(ctx, db, data, persistingData); err != nil {
				return apperrors.Wrap(err)
			}
		}

		err = persistSettingUpdate(ctx, db, data, persistingData)
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

	return &UpdateSettingResp{}, nil
}

func loadSettingForUpdate(
	ctx context.Context,
	db database.IDB,
	req *UpdateSettingReq,
	data *UpdateSettingData,
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
	if req.UpdateVer != setting.UpdateVer {
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}
	data.Setting = setting

	// If name changes, validate the new one
	if data.VerifyingName != "" && !strings.EqualFold(setting.Name, data.VerifyingName) {
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

func prepareSettingUpdate(
	_ *UpdateSettingReq,
	data *UpdateSettingData,
	persistingData *PersistingSettingData,
) {
	timeNow := timeutil.NowUTC()
	setting := data.Setting
	setting.UpdateVer++
	setting.UpdatedAt = timeNow

	persistingData.Setting = setting
}

func persistSettingUpdate(
	ctx context.Context,
	db database.IDB,
	data *UpdateSettingData,
	persistingData *PersistingSettingData,
) error {
	err := data.SettingRepo.Update(ctx, db, persistingData.Setting)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
