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
	"github.com/localpaas/localpaas/localpaas_app/pkg/entityutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/strutil"
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

func (uc *BaseSettingUC) loadSettingByID(
	ctx context.Context,
	db database.IDB,
	req *BaseSettingReq,
	id string,
	requireActive bool,
	opts ...bunex.SelectQueryOption,
) (setting *entity.Setting, err error) {
	loadOpts := append([]bunex.SelectQueryOption{}, opts...)
	switch req.Scope {
	case base.SettingScopeGlobal:
		loadOpts = append(loadOpts, bunex.SelectWhere("setting.object_id IS NULL"))
		setting, err = uc.SettingRepo.GetByID(ctx, db, req.Type, id, requireActive, loadOpts...)
	case base.SettingScopeProject:
		setting, err = uc.SettingRepo.GetByIDAndProject(ctx, db, req.Type, id, req.ObjectID,
			requireActive, loadOpts...)
	case base.SettingScopeApp:
		setting, err = uc.SettingRepo.GetByIDAndApp(ctx, db, req.Type, id, req.ParentObjectID, req.ObjectID,
			requireActive, loadOpts...)
	case base.SettingScopeUser:
		setting, err = uc.SettingRepo.GetByIDAndUser(ctx, db, req.Type, id, req.ObjectID,
			requireActive, loadOpts...)
	case base.SettingScopeNone:
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return setting, nil
}

func (uc *BaseSettingUC) checkNameConflict(
	ctx context.Context,
	db database.IDB,
	req *BaseSettingReq,
	name string,
) (err error) {
	if name == "" {
		return nil
	}
	var setting *entity.Setting
	switch req.Scope {
	case base.SettingScopeGlobal:
		setting, err = uc.SettingRepo.GetByName(ctx, db, req.Type, name, false,
			bunex.SelectWhere("setting.object_id IS NULL"),
		)
	case base.SettingScopeProject:
		setting, err = uc.SettingRepo.GetByNameAndProject(ctx, db, req.Type, name, req.ObjectID, false)
	case base.SettingScopeApp:
		setting, err = uc.SettingRepo.GetByNameAndApp(ctx, db, req.Type, name, req.ParentObjectID, req.ObjectID, false)
	case base.SettingScopeUser:
		setting, err = uc.SettingRepo.GetByNameAndUser(ctx, db, req.Type, name, req.ObjectID, false)
	case base.SettingScopeNone:
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

func (uc *BaseSettingUC) checkRefSettingsExistence(
	ctx context.Context,
	db database.IDB,
	req *BaseSettingReq,
	refSettingIDs []string,
	requireActive bool,
) (err error) {
	if len(refSettingIDs) == 0 {
		return nil
	}
	listOpts := []bunex.SelectQueryOption{
		bunex.SelectWhere("setting.id IN (?)", bunex.In(refSettingIDs)),
	}
	if requireActive {
		listOpts = append(listOpts, bunex.SelectWhere("setting.status = ?", base.SettingStatusActive))
	}
	var settings []*entity.Setting
	switch req.Scope {
	case base.SettingScopeGlobal:
		listOpts = append(listOpts, bunex.SelectWhere("setting.object_id IS NULL"))
		settings, _, err = uc.SettingRepo.List(ctx, db, nil, listOpts...)
	case base.SettingScopeProject:
		settings, _, err = uc.SettingRepo.ListByProject(ctx, db, req.ObjectID, nil, listOpts...)
	case base.SettingScopeApp:
		settings, _, err = uc.SettingRepo.ListByApp(ctx, db, req.ParentObjectID, req.ObjectID, nil, listOpts...)
	case base.SettingScopeUser:
		settings, _, err = uc.SettingRepo.ListByUser(ctx, db, req.ObjectID, nil, listOpts...)
	case base.SettingScopeNone:
	}
	if err != nil {
		return apperrors.Wrap(err)
	}
	for _, refSettingID := range refSettingIDs {
		found := entityutil.FindByID(settings, refSettingID)
		if found == nil {
			return apperrors.NewNotFound("Setting").WithMsgLog("setting %s not found", refSettingID)
		}
	}
	return nil
}

func (uc *BaseSettingUC) ensureSettingDefaultUniqueness(
	ctx context.Context,
	db database.IDB,
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
	case base.SettingScopeNone:
	}

	err := uc.SettingRepo.UpdateClearDefaultFlag(ctx, db, req.Type, setting.ID, opts...)
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
