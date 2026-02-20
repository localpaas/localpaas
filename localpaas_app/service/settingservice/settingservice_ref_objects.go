package settingservice

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/entityutil"
)

func (s *settingService) LoadReferenceObjects(
	ctx context.Context,
	db database.IDB,
	scope base.SettingScope,
	objectID string,
	parentObjectID string,
	requireActive bool,
	errorIfUnavail bool,
	inSettings ...*entity.Setting,
) (refObjects *entity.RefObjects, err error) {
	allRefIDs := &entity.RefObjectIDs{}
	for _, setting := range inSettings {
		allRefIDs.AddRefIDs(setting.MustGetRefObjectIDs())
	}
	return s.LoadReferenceObjectsByIDs(ctx, db, scope, objectID, parentObjectID, requireActive,
		errorIfUnavail, allRefIDs)
}

func (s *settingService) LoadReferenceObjectsByIDs(
	ctx context.Context,
	db database.IDB,
	scope base.SettingScope,
	objectID string,
	parentObjectID string,
	requireActive bool,
	errorIfUnavail bool,
	refIDs *entity.RefObjectIDs,
) (refObjects *entity.RefObjects, err error) {
	refObjects = &entity.RefObjects{}

	if refIDs == nil || !refIDs.HasData() {
		return refObjects, nil
	}

	// Load ref users
	if len(refIDs.RefUserIDs) > 0 {
		refObjects.RefUsers, err = s.userService.LoadUsers(ctx, db, refIDs.RefUserIDs, errorIfUnavail)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
	}

	// Make sure the current app id is in the list
	if scope == base.SettingScopeApp && !gofn.Contain(refIDs.RefAppIDs, objectID) {
		refIDs.RefAppIDs = append(refIDs.RefAppIDs, objectID)
	}
	// Load ref apps
	if len(refIDs.RefAppIDs) > 0 {
		var projectID string
		switch scope { //nolint:exhaustive
		case base.SettingScopeProject:
			projectID = objectID
		case base.SettingScopeApp:
			projectID = parentObjectID
		}
		refObjects.RefApps, err = s.LoadReferenceApps(ctx, db, projectID, requireActive,
			errorIfUnavail, refIDs.RefAppIDs)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
	}

	// Load ref settings
	if len(refIDs.RefSettingIDs) > 0 {
		if scope == base.SettingScopeApp && parentObjectID == "" {
			app := refObjects.RefApps[objectID]
			if app != nil && app.Project != nil {
				parentObjectID = app.Project.ID
			}
		}
		refObjects.RefSettings, err = s.LoadReferenceSettings(ctx, db, scope, objectID, parentObjectID,
			requireActive, errorIfUnavail, refIDs.RefSettingIDs)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
	}

	// Calculate recursive ref IDs to load
	newRecursiveRefIDs := refIDs.GetRecursiveRefObjectIDs(refObjects)
	if !newRecursiveRefIDs.HasData() {
		return refObjects, nil
	}

	newRecursiveRefObjects, err := s.LoadReferenceObjectsByIDs(ctx, db, scope, objectID, parentObjectID,
		requireActive, errorIfUnavail, newRecursiveRefIDs)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	refObjects.AddRefObjects(newRecursiveRefObjects)

	return refObjects, nil
}

func (s *settingService) LoadReferenceSettings(
	ctx context.Context,
	db database.IDB,
	scope base.SettingScope,
	objectID string,
	parentObjectID string,
	requireActive bool,
	errorIfUnavail bool,
	settingIDs []string,
) (settingMap map[string]*entity.Setting, err error) {
	settingIDs = gofn.ToSet(settingIDs)
	listOpts := []bunex.SelectQueryOption{
		bunex.SelectWhereIn("setting.id IN (?)", settingIDs...),
	}
	if requireActive {
		listOpts = append(listOpts, bunex.SelectWhere("setting.status = ?", base.SettingStatusActive))
	}

	var settings []*entity.Setting
	switch scope {
	case base.SettingScopeNone:
		settings, _, err = s.settingRepo.List(ctx, db, nil, listOpts...)
	case base.SettingScopeGlobal:
		settings, _, err = s.settingRepo.ListGlobally(ctx, db, nil, listOpts...)
	case base.SettingScopeProject:
		settings, _, err = s.settingRepo.ListByProject(ctx, db, objectID, nil, listOpts...)
	case base.SettingScopeApp:
		settings, _, err = s.settingRepo.ListByApp(ctx, db, objectID, parentObjectID, nil, listOpts...)
	case base.SettingScopeUser:
		settings, _, err = s.settingRepo.ListByUser(ctx, db, objectID, nil, listOpts...)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	settingMap = entityutil.SliceToIDMap(settings)

	// Check setting availability
	if errorIfUnavail {
		for _, id := range settingIDs {
			if _, exists := settingMap[id]; !exists {
				return nil, apperrors.NewNotFound("Setting").
					WithMsgLog("setting %s not found or expired", id)
			}
		}
	}

	return settingMap, nil
}

func (s *settingService) LoadReferenceApps(
	ctx context.Context,
	db database.IDB,
	projectID string,
	requireActive bool,
	errorIfUnavail bool,
	appIDs []string,
) (appMap map[string]*entity.App, err error) {
	appIDs = gofn.ToSet(appIDs)
	opts := []bunex.SelectQueryOption{
		bunex.SelectExcludeColumns(entity.AppDefaultExcludeColumns...),
		bunex.SelectRelation("Project",
			bunex.SelectExcludeColumns(entity.ProjectDefaultExcludeColumns...),
		),
	}
	if projectID != "" {
		opts = append(opts, bunex.SelectWhere("app.project_id = ?", projectID))
	}
	if requireActive {
		opts = append(opts, bunex.SelectWhere("app.status = ?", base.AppStatusActive))
	}

	apps, err := s.appRepo.ListByIDs(ctx, db, "", appIDs, opts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	appMap = entityutil.SliceToIDMap(apps)

	for _, id := range appIDs {
		app, exists := appMap[id]
		if errorIfUnavail && !exists {
			return nil, apperrors.NewNotFound("App").
				WithMsgLog("app %s not found or inactive", id)
		}
		if requireActive && app.Project != nil && app.Project.Status != base.ProjectStatusActive {
			app.Project = nil
		}
		if errorIfUnavail && app.Project == nil {
			return nil, apperrors.NewNotFound("Project").
				WithMsgLog("project %s not found", app.ProjectID)
		}
	}

	return appMap, nil
}
