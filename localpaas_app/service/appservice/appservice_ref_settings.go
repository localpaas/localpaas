package appservice

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/strutil"
)

func (s *appService) LoadReferenceSettings(
	ctx context.Context,
	db database.IDB,
	app *entity.App,
	appSettings ...*entity.Setting,
) (settingMap map[string]*entity.Setting, err error) {
	settingIDMap := make(map[string]base.SettingType, 10) //nolint

	for _, setting := range appSettings {
		switch setting.Type { //nolint:exhaustive
		case base.SettingTypeAppHttp:
			httpSettings, err := setting.AsAppHttpSettings()
			if err != nil {
				return nil, apperrors.Wrap(err)
			}
			for _, domain := range httpSettings.Domains {
				settingIDMap[domain.SSLCert.ID] = base.SettingTypeSSL
				settingIDMap[domain.BasicAuth.ID] = base.SettingTypeBasicAuth
			}

		case base.SettingTypeAppDeployment:
			deplSettings, err := setting.AsAppDeploymentSettings()
			if err != nil {
				return nil, apperrors.Wrap(err)
			}
			if deplSettings.ImageSource != nil {
				settingIDMap[deplSettings.ImageSource.RegistryAuth.ID] = base.SettingTypeRegistryAuth
			}
			if deplSettings.RepoSource != nil {
				settingIDMap[deplSettings.RepoSource.Credentials.ID] = deplSettings.RepoSource.Credentials.Type
				settingIDMap[deplSettings.RepoSource.RegistryAuth.ID] = base.SettingTypeRegistryAuth
			}
			// TODO: handle deplSettings.TarballSource

		default:
			return nil, apperrors.New(apperrors.ErrUnsupported).
				WithMsgLog("unsupported app setting type: %s", setting.Type)
		}
	}

	// Remove empty key
	delete(settingIDMap, "")

	settings, _, err := s.settingRepo.ListByApp(ctx, db, app.ProjectID, app.ID, nil,
		bunex.SelectWhere("setting.id IN (?)", bunex.In(gofn.MapKeys(settingIDMap))),
		bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	settingMap = make(map[string]*entity.Setting, len(settings))
	for _, setting := range settings {
		settingMap[setting.ID] = setting
	}

	// Check setting existence
	for id, settingType := range settingIDMap {
		if _, exists := settingMap[id]; !exists {
			return nil, apperrors.NewNotFound(strutil.ToPascalCase(string(settingType))).
				WithMsgLog("%s %s not found or expired", settingType, id)
		}
	}

	return settingMap, nil
}
