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
	"github.com/localpaas/localpaas/localpaas_app/repository"
)

type GetSettingReq struct {
	ID        string            `json:"-" mapstructure:"-"`
	Type      base.SettingType  `json:"-" mapstructure:"-"`
	Scope     base.SettingScope `json:"-" mapstructure:"-"`
	ProjectID string            `json:"-" mapstructure:"-"`
	AppID     string            `json:"-" mapstructure:"-"`
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
	loadOpts := []bunex.SelectQueryOption{
		bunex.SelectWhere("setting.type = ?", req.Type),
		bunex.SelectWhereIf(req.Scope == base.SettingScopeGlobal, "setting.object_id IS NULL"),
	}
	loadOpts = append(loadOpts, data.ExtraLoadOpts...)

	setting, err := data.SettingRepo.GetByIDEx(ctx, db, req.Type, req.ProjectID, req.AppID, req.ID,
		false, loadOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return setting, nil
}
