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

type ListSettingReq struct {
	BaseSettingReq
	Status []base.SettingStatus `json:"-" mapstructure:"status"`
	Search string               `json:"-" mapstructure:"search"`

	Paging basedto.Paging `json:"-"`
}

func (req *ListSettingReq) Validate() (validators []vld.Validator) {
	validators = append(validators, basedto.ValidateSlice(req.Status, true, 0, base.AllSettingStatuses, "status")...)
	return
}

type ListSettingResp struct {
	Meta *basedto.Meta     `json:"meta"`
	Data []*entity.Setting `json:"data"`
}

type ListSettingData struct {
	SettingRepo   repository.SettingRepo
	ExtraLoadOpts []bunex.SelectQueryOption
}

func ListSetting(
	ctx context.Context,
	db database.IDB,
	auth *basedto.Auth,
	req *ListSettingReq,
	data *ListSettingData,
) (_ *ListSettingResp, err error) {
	listOpts := []bunex.SelectQueryOption{
		bunex.SelectWhere("setting.type = ?", req.Type),
	}
	if len(req.Status) > 0 {
		listOpts = append(listOpts, bunex.SelectWhere("setting.status IN (?)", bunex.In(req.Status)))
	}
	if req.Search != "" {
		keyword := bunex.MakeLikeOpStr(req.Search, true)
		listOpts = append(listOpts,
			bunex.SelectWhereGroup(
				bunex.SelectWhere("setting.name ILIKE ?", keyword),
			),
		)
	}
	if len(auth.AllowObjectIDs) > 0 {
		listOpts = append(listOpts,
			bunex.SelectWhere("setting.id IN (?)", bunex.In(auth.AllowObjectIDs)),
		)
	}
	listOpts = append(listOpts, data.ExtraLoadOpts...)

	var settings []*entity.Setting
	var paging *basedto.PagingMeta

	switch req.Scope {
	case base.SettingScopeGlobal:
		listOpts = append(listOpts, bunex.SelectWhere("setting.object_id IS NULL"))
		settings, paging, err = data.SettingRepo.List(ctx, db, &req.Paging, listOpts...)
	case base.SettingScopeProject:
		settings, paging, err = data.SettingRepo.ListByProject(ctx, db, req.ObjectID,
			&req.Paging, listOpts...)
	case base.SettingScopeApp:
		settings, paging, err = data.SettingRepo.ListByApp(ctx, db, req.ParentObjectID, req.ObjectID,
			&req.Paging, listOpts...)
	case base.SettingScopeUser:
		settings, paging, err = data.SettingRepo.ListByUser(ctx, db, req.ObjectID,
			&req.Paging, listOpts...)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	for _, setting := range settings {
		setting.CurrentObjectID = req.ObjectID
	}

	return &ListSettingResp{
		Meta: &basedto.Meta{Page: paging},
		Data: settings,
	}, nil
}
