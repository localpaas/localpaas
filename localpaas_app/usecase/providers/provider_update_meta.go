package providers

import (
	"context"
	"time"

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

type UpdateSettingMetaReq struct {
	ID         string              `json:"-"`
	Type       base.SettingType    `json:"-"`
	ProjectID  string              `json:"-"`
	AppID      string              `json:"-"`
	GlobalOnly bool                `json:"-"`
	Status     *base.SettingStatus `json:"status"`
	ExpireAt   *time.Time          `json:"expireAt"`
	UpdateVer  int                 `json:"updateVer"`
}

type UpdateSettingMetaResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}

type UpdateSettingMetaData struct {
	Setting *entity.Setting

	SettingRepo   repository.SettingRepo
	ExtraLoadOpts []bunex.SelectQueryOption

	AfterLoading     func(context.Context, database.Tx, *UpdateSettingMetaData) error
	BeforePersisting func(context.Context, database.Tx, *UpdateSettingMetaData, *PersistingSettingMetaData) error
	AfterPersisting  func(context.Context, database.Tx, *UpdateSettingMetaData, *PersistingSettingMetaData) error
}

type PersistingSettingMetaData struct {
	Setting *entity.Setting
}

func UpdateSettingMeta(
	ctx context.Context,
	db database.IDB,
	req *UpdateSettingMetaReq,
	data *UpdateSettingMetaData,
) (*UpdateSettingMetaResp, error) {
	err := transaction.Execute(ctx, db, func(db database.Tx) error {
		err := loadSettingForUpdateMeta(ctx, db, req, data)
		if err != nil {
			return apperrors.Wrap(err)
		}

		if data.AfterLoading != nil {
			if err := data.AfterLoading(ctx, db, data); err != nil {
				return apperrors.Wrap(err)
			}
		}

		persistingData := &PersistingSettingMetaData{}
		prepareSettingMetaUpdate(req, data, persistingData)
		if data.BeforePersisting != nil {
			if err := data.BeforePersisting(ctx, db, data, persistingData); err != nil {
				return apperrors.Wrap(err)
			}
		}

		err = persistSettingMetaUpdate(ctx, db, data, persistingData)
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

	return &UpdateSettingMetaResp{}, nil
}

func loadSettingForUpdateMeta(
	ctx context.Context,
	db database.IDB,
	req *UpdateSettingMetaReq,
	data *UpdateSettingMetaData,
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

	return nil
}

func prepareSettingMetaUpdate(
	req *UpdateSettingMetaReq,
	data *UpdateSettingMetaData,
	persistingData *PersistingSettingMetaData,
) {
	timeNow := timeutil.NowUTC()
	setting := data.Setting
	setting.UpdateVer++
	setting.UpdatedAt = timeNow

	if req.Status != nil {
		setting.Status = *req.Status
	}
	if req.ExpireAt != nil {
		setting.ExpireAt = *req.ExpireAt
	}

	persistingData.Setting = setting
}

func persistSettingMetaUpdate(
	ctx context.Context,
	db database.IDB,
	data *UpdateSettingMetaData,
	persistingData *PersistingSettingMetaData,
) error {
	err := data.SettingRepo.Update(ctx, db, persistingData.Setting,
		bunex.UpdateColumns("update_ver", "updated_at", "status", "expire_at"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
