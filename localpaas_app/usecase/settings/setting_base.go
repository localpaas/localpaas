package settings

import (
	"context"
	"errors"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/pkg/strutil"
	"github.com/localpaas/localpaas/localpaas_app/repository"
)

type BaseSettingReq struct {
	Type           base.SettingType  `json:"-" mapstructure:"-"`
	Scope          base.SettingScope `json:"-" mapstructure:"-"`
	ObjectID       string            `json:"-" mapstructure:"-"`
	ParentObjectID string            `json:"-" mapstructure:"-"`
}

type BaseSettingResp struct {
	ID              string             `json:"id"`
	Type            base.SettingType   `json:"type"`
	Name            string             `json:"name"`
	Kind            string             `json:"kind,omitempty"`
	Status          base.SettingStatus `json:"status"`
	Inherited       bool               `json:"inherited,omitempty"`
	AvailInProjects bool               `json:"availableInProjects,omitempty"`
	Default         bool               `json:"default,omitempty"`
	UpdateVer       int                `json:"updateVer"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ExpireAt  *time.Time `json:"expireAt,omitempty" copy:",nilonzero"`
}

func loadSettingByID(
	ctx context.Context,
	db database.IDB,
	settingRepo repository.SettingRepo,
	req *BaseSettingReq,
	id string,
	requireActive bool, //nolint:unparam
	opts ...bunex.SelectQueryOption,
) (setting *entity.Setting, err error) {
	loadOpts := []bunex.SelectQueryOption{
		bunex.SelectWhere("setting.type = ?", req.Type),
	}
	loadOpts = append(loadOpts, opts...)

	switch req.Scope {
	case base.SettingScopeGlobal:
		loadOpts = append(loadOpts, bunex.SelectWhere("setting.object_id IS NULL"))
		setting, err = settingRepo.GetByID(ctx, db, req.Type, id, requireActive, loadOpts...)
	case base.SettingScopeProject:
		setting, err = settingRepo.GetByIDAndProject(ctx, db, req.Type, id, req.ObjectID,
			requireActive, loadOpts...)
	case base.SettingScopeApp:
		setting, err = settingRepo.GetByIDAndApp(ctx, db, req.Type, id, req.ParentObjectID, req.ObjectID,
			requireActive, loadOpts...)
	case base.SettingScopeUser:
		setting, err = settingRepo.GetByIDAndUser(ctx, db, req.Type, id, req.ObjectID,
			requireActive, loadOpts...)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return setting, nil
}

func checkNameConflict(
	ctx context.Context,
	db database.IDB,
	settingRepo repository.SettingRepo,
	req *BaseSettingReq,
	name string,
) (err error) {
	if name == "" {
		return nil
	}
	var setting *entity.Setting
	switch req.Scope {
	case base.SettingScopeGlobal:
		setting, err = settingRepo.GetByName(ctx, db, req.Type, name, false,
			bunex.SelectWhere("setting.object_id IS NULL"),
		)
	case base.SettingScopeProject:
		setting, err = settingRepo.GetByNameAndProject(ctx, db, req.Type, name, req.ObjectID, false)
	case base.SettingScopeApp:
		setting, err = settingRepo.GetByNameAndApp(ctx, db, req.Type, name, req.ParentObjectID, req.ObjectID, false)
	case base.SettingScopeUser:
		setting, err = settingRepo.GetByNameAndUser(ctx, db, req.Type, name, req.ObjectID, false)
	}
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if setting != nil {
		return apperrors.NewAlreadyExist(strutil.ToPascalCase(string(req.Type))).
			WithMsgLog("%s '%s' already exists", req.Type, setting.Name)
	}
	return nil
}

func ensureSettingDefaultUniqueness(
	ctx context.Context,
	db database.IDB,
	settingRepo repository.SettingRepo,
	req *BaseSettingReq,
	setting *entity.Setting,
) error {
	opts := []bunex.UpdateQueryOption{
		bunex.UpdateWithDeleted(),
	}
	switch req.Scope {
	case base.SettingScopeGlobal:
		opts = append(opts, bunex.UpdateWhere("setting.object_id IS NULL"))
	case base.SettingScopeProject:
		opts = append(opts, bunex.UpdateWhere("setting.object_id = ?", req.ObjectID))
	case base.SettingScopeApp:
		opts = append(opts, bunex.UpdateWhere("setting.object_id = ?", req.ObjectID))
	case base.SettingScopeUser:
		opts = append(opts, bunex.UpdateWhere("setting.object_id = ?", req.ObjectID))
	}

	err := settingRepo.UpdateClearDefaultFlag(ctx, db, req.Type, setting.ID, opts...)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func TransformSettingBase(setting *entity.Setting) (resp *BaseSettingResp, err error) {
	if setting == nil {
		return nil, nil
	}
	if err = copier.Copy(&resp, setting); err != nil {
		return nil, apperrors.Wrap(err)
	}
	if setting.ObjectID != setting.CurrentObjectID {
		resp.Inherited = true
	}
	return resp, nil
}
