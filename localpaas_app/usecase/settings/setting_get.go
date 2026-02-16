package settings

import (
	"context"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
)

type GetSettingReq struct {
	BaseSettingReq
	ID string `json:"-" mapstructure:"-"`
}

func (req *GetSettingReq) Validate() (validators []vld.Validator) {
	return
}

type GetSettingResp struct {
	Data       *entity.Setting
	RefObjects *entity.RefObjects
}

type GetSettingData struct {
	ExtraLoadOpts []bunex.SelectQueryOption
}

func (uc *BaseSettingUC) GetSetting(
	ctx context.Context,
	auth *basedto.Auth,
	req *GetSettingReq,
	data *GetSettingData,
) (*GetSettingResp, error) {
	setting, err := uc.loadSettingByID(ctx, uc.DB, &req.BaseSettingReq, req.ID,
		false, data.ExtraLoadOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if setting != nil {
		setting.CurrentObjectID = req.ObjectID
	}

	refObjects := &entity.RefObjects{}
	err = uc.loadRefObjects(ctx, uc.DB, &req.BaseSettingReq, []*entity.Setting{setting}, refObjects)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &GetSettingResp{
		Data:       setting,
		RefObjects: refObjects,
	}, nil
}

func (uc *BaseSettingUC) GetSettingByID(
	ctx context.Context,
	db database.IDB,
	req *BaseSettingReq,
	id string,
	requireActive bool,
	extraLoadOpts ...bunex.SelectQueryOption,
) (*entity.Setting, error) {
	setting, err := uc.loadSettingByID(ctx, db, req, id, requireActive, extraLoadOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if setting != nil {
		setting.CurrentObjectID = req.ObjectID
	}
	return setting, nil
}
