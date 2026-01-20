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
	Type      base.SettingType     `json:"-" mapstructure:"-"`
	Scope     base.SettingScope    `json:"-" mapstructure:"-"`
	ProjectID string               `json:"-" mapstructure:"-"`
	AppID     string               `json:"-" mapstructure:"-"`
	Status    []base.SettingStatus `json:"-" mapstructure:"status"`
	Search    string               `json:"-" mapstructure:"search"`

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
) (*ListSettingResp, error) {
	listOpts := []bunex.SelectQueryOption{
		bunex.SelectWhere("setting.type = ?", req.Type),
		bunex.SelectWhereIf(req.Scope == base.SettingScopeGlobal, "setting.object_id IS NULL"),
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

	settings, paging, err := data.SettingRepo.List(ctx, db, req.ProjectID, req.AppID, &req.Paging, listOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &ListSettingResp{
		Meta: &basedto.Meta{Page: paging},
		Data: settings,
	}, nil
}
