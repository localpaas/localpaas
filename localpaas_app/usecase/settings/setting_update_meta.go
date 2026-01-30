package settings

import (
	"context"
	"time"

	"github.com/tiendc/gofn"

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
	BaseSettingReq
	ID                  string              `json:"-"`
	Status              *base.SettingStatus `json:"status"`
	ExpireAt            *time.Time          `json:"expireAt"`
	AvailableInProjects *bool               `json:"availableInProjects"`
	Default             *bool               `json:"default"`
	UpdateVer           int                 `json:"updateVer"`
}

type UpdateSettingMetaResp struct {
	Meta *basedto.Meta `json:"meta"`
}

type UpdateSettingMetaData struct {
	Setting *entity.Setting

	SettingRepo       repository.SettingRepo
	DefaultMustUnique bool
	ExtraLoadOpts     []bunex.SelectQueryOption

	AfterLoading     func(context.Context, database.Tx, *UpdateSettingMetaData) error
	BeforePersisting func(context.Context, database.Tx, *UpdateSettingMetaData, *PersistingSettingMetaData) error
	AfterPersisting  func(context.Context, database.Tx, *UpdateSettingMetaData, *PersistingSettingMetaData) error

	oldDefaultFlag bool
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

		err = persistSettingMetaUpdate(ctx, db, req, data, persistingData)
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
	}
	loadOpts = append(loadOpts, data.ExtraLoadOpts...)

	setting, err := loadSettingByID(ctx, db, data.SettingRepo, &req.BaseSettingReq, req.ID,
		false, false, loadOpts...)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if req.UpdateVer != setting.UpdateVer {
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}
	data.Setting = setting

	if req.Scope != base.SettingScopeGlobal && setting.ObjectID != req.ObjectID {
		return apperrors.New(apperrors.ErrOwnSettingRequired).
			WithMsgLog("imported or inherited setting is not allowed to update")
	}

	data.oldDefaultFlag = setting.Default
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
	if req.AvailableInProjects != nil {
		setting.AvailInProjects = gofn.If(req.Scope != base.SettingScopeGlobal, false, *req.AvailableInProjects)
	}
	if req.Default != nil {
		setting.Default = *req.Default
	}

	persistingData.Setting = setting
}

func persistSettingMetaUpdate(
	ctx context.Context,
	db database.IDB,
	req *UpdateSettingMetaReq,
	data *UpdateSettingMetaData,
	persistingData *PersistingSettingMetaData,
) error {
	err := data.SettingRepo.Update(ctx, db, persistingData.Setting,
		bunex.UpdateColumns("update_ver", "updated_at", "status", "expire_at", "avail_in_projects", "is_default"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}

	if data.DefaultMustUnique && !data.oldDefaultFlag && persistingData.Setting.Default {
		if data.DefaultMustUnique && persistingData.Setting.Default {
			err = ensureSettingDefaultUniqueness(ctx, db, data.SettingRepo, &req.BaseSettingReq, persistingData.Setting)
			if err != nil {
				return apperrors.Wrap(err)
			}
		}
	}

	return nil
}
