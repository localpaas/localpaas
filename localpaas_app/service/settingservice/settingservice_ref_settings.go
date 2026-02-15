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

func (s *settingService) LoadReferenceSettings(
	ctx context.Context,
	db database.IDB,
	project *entity.Project,
	app *entity.App,
	requireActive bool,
	settingIDs []string,
) (settingMap map[string]*entity.Setting, err error) {
	settingIDs = gofn.ToSet(settingIDs)
	opts := []bunex.SelectQueryOption{
		bunex.SelectWhereIn("setting.id IN (?)", settingIDs...),
	}
	if requireActive {
		opts = append(opts, bunex.SelectWhere("setting.status = ?", base.SettingStatusActive))
	}

	var settings []*entity.Setting
	switch {
	case app != nil:
		settings, _, err = s.settingRepo.ListByAppObject(ctx, db, app, nil, opts...)
	case project != nil:
		settings, _, err = s.settingRepo.ListByProject(ctx, db, project.ID, nil, opts...)
	default:
		settings, _, err = s.settingRepo.ListGlobally(ctx, db, nil, opts...)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	settingMap = entityutil.SliceToIDMap(settings)

	// Check setting existence
	for _, id := range settingIDs {
		if _, exists := settingMap[id]; !exists {
			return nil, apperrors.NewNotFound("Setting").
				WithMsgLog("setting %s not found or expired", id)
		}
	}

	return settingMap, nil
}

func (s *settingService) LoadReferenceSettingsFor(
	ctx context.Context,
	db database.IDB,
	project *entity.Project,
	app *entity.App,
	requireActive bool,
	inSettings ...*entity.Setting,
) (settingMap map[string]*entity.Setting, err error) {
	settingIDMap := make(map[string]struct{}, 10) //nolint
	for _, setting := range inSettings {
		for _, settingID := range setting.MustGetRefSettingIDs() {
			settingIDMap[settingID] = struct{}{}
		}
	}
	return s.LoadReferenceSettings(ctx, db, project, app, requireActive, gofn.MapKeys(settingIDMap))
}
