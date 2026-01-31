package appservice

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

func (s *appService) LoadSettings(
	ctx context.Context,
	db database.IDB,
	app *entity.App,
	settingIDs []string,
	requireActive bool,
) (settingMap map[string]*entity.Setting, err error) {
	opts := []bunex.SelectQueryOption{
		bunex.SelectWhere("setting.id IN (?)", bunex.In(settingIDs)),
	}
	if requireActive {
		opts = append(opts, bunex.SelectWhere("setting.status = ?", base.SettingStatusActive))
	}
	settings, _, err := s.settingRepo.ListByAppObject(ctx, db, app, nil, opts...)
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

func (s *appService) LoadReferenceSettings(
	ctx context.Context,
	db database.IDB,
	app *entity.App,
	requireActive bool,
	appSettings ...*entity.Setting,
) (settingMap map[string]*entity.Setting, err error) {
	settingIDMap := make(map[string]struct{}, 10) //nolint
	for _, setting := range appSettings {
		for _, settingID := range setting.RefIDs {
			settingIDMap[settingID] = struct{}{}
		}
	}

	opts := []bunex.SelectQueryOption{
		bunex.SelectWhere("setting.id IN (?)", bunex.In(gofn.MapKeys(settingIDMap))),
	}
	if requireActive {
		opts = append(opts, bunex.SelectWhere("setting.status = ?", base.SettingStatusActive))
	}
	settings, _, err := s.settingRepo.ListByAppObject(ctx, db, app, nil, opts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	settingMap = entityutil.SliceToIDMap(settings)

	// Check setting existence
	for id := range settingIDMap {
		if _, exists := settingMap[id]; !exists {
			return nil, apperrors.NewNotFound("Setting").
				WithMsgLog("setting %s not found or expired", id)
		}
	}

	return settingMap, nil
}
