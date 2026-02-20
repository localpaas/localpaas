package settings

import (
	"context"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
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
	Meta       *basedto.ListMeta
	Data       []*entity.Setting
	RefObjects *entity.RefObjects
}

type ListSettingData struct {
	ExtraLoadOpts []bunex.SelectQueryOption
}

func (uc *BaseSettingUC) ListSetting(
	ctx context.Context,
	auth *basedto.Auth,
	req *ListSettingReq,
	data *ListSettingData,
) (_ *ListSettingResp, err error) {
	db := uc.DB
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
		settings, paging, err = uc.SettingRepo.ListGlobally(ctx, db, &req.Paging, listOpts...)
	case base.SettingScopeProject:
		settings, paging, err = uc.SettingRepo.ListByProject(ctx, db, req.ObjectID,
			&req.Paging, listOpts...)
	case base.SettingScopeApp:
		settings, paging, err = uc.SettingRepo.ListByApp(ctx, db, req.ObjectID, req.ParentObjectID,
			&req.Paging, listOpts...)
	case base.SettingScopeUser:
		settings, paging, err = uc.SettingRepo.ListByUser(ctx, db, req.ObjectID,
			&req.Paging, listOpts...)
	case base.SettingScopeNone:
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	for _, setting := range settings {
		setting.CurrentObjectID = req.ObjectID
	}

	refObjects, err := uc.SettingService.LoadReferenceObjects(ctx, uc.DB, req.Scope, req.ObjectID,
		req.ParentObjectID, true, false, settings...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &ListSettingResp{
		Meta:       &basedto.ListMeta{Page: paging},
		Data:       settings,
		RefObjects: refObjects,
	}, nil
}
