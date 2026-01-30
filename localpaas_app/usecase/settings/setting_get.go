package settings

import (
	"context"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/repository"
)

type GetSettingReq struct {
	BaseSettingReq
	ID string `json:"-" mapstructure:"-"`
}

func (req *GetSettingReq) Validate() (validators []vld.Validator) {
	return
}

type GetSettingData struct {
	SettingRepo   repository.SettingRepo
	ExtraLoadOpts []bunex.SelectQueryOption
}

func GetSetting(
	ctx context.Context,
	db database.IDB,
	auth *basedto.Auth,
	req *GetSettingReq,
	data *GetSettingData,
) (*entity.Setting, error) {
	setting, err := loadSettingByID(ctx, db, data.SettingRepo, &req.BaseSettingReq, req.ID,
		false, true, data.ExtraLoadOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if setting != nil {
		setting.CurrentObjectID = req.ObjectID
	}
	return setting, nil
}

func GetSettingByID(
	ctx context.Context,
	db database.IDB,
	settingRepo repository.SettingRepo,
	req *BaseSettingReq,
	id string,
	requireActive bool,
	loadRefSettings bool,
	extraLoadOpts ...bunex.SelectQueryOption,
) (*entity.Setting, error) {
	setting, err := loadSettingByID(ctx, db, settingRepo, req, id, requireActive,
		loadRefSettings, extraLoadOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if setting != nil {
		setting.CurrentObjectID = req.ObjectID
	}
	return setting, nil
}
