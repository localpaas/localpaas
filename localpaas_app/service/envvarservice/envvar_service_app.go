package envvarservice

import (
	"context"
	"fmt"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
)

//nolint:gocognit
func (s *envVarService) BuildAppEnv(ctx context.Context, db database.IDB, app *entity.App, buildPhase bool) (
	res []string, err error) {
	objectIDs := gofn.ToSliceSkippingZero(app.ID, app.ParentID, app.ProjectID)
	settings, _, err := s.settingRepo.List(ctx, db, nil,
		bunex.SelectWhere("setting.type IN (?)",
			bunex.In([]base.SettingType{base.SettingTypeEnvVar, base.SettingTypeSecret})),
		bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
		bunex.SelectWhere("setting.object_id IN (?)", bunex.In(objectIDs)),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	mapAppEnv := make(map[string]string, 20)                   //nolint
	mapParentAppEnv := make(map[string]string, 20)             //nolint
	mapProjectEnv := make(map[string]string, 20)               //nolint
	mapAppSecret := make(map[string]*entity.Setting, 10)       //nolint
	mapParentAppSecret := make(map[string]*entity.Setting, 10) //nolint
	mapProjectSecret := make(map[string]*entity.Setting, 10)   //nolint

	for _, setting := range settings {
		if setting.Type == base.SettingTypeEnvVar {
			vars, err := setting.ParseEnvVars()
			if err != nil {
				return nil, apperrors.Wrap(err)
			}
			switch setting.ObjectID {
			case app.ID:
				for _, env := range vars.Data {
					if env.IsBuildEnv == buildPhase {
						mapAppEnv[env.Key] = env.Value
					}
				}
			case app.ParentID:
				for _, env := range vars.Data {
					if env.IsBuildEnv == buildPhase {
						mapParentAppEnv[env.Key] = env.Value
					}
				}
			case app.ProjectID:
				for _, env := range vars.Data {
					if env.IsBuildEnv == buildPhase {
						mapProjectEnv[env.Key] = env.Value
					}
				}
			}
			continue
		}

		if setting.Type == base.SettingTypeSecret {
			switch setting.ObjectID {
			case app.ID:
				mapAppSecret[setting.Name] = setting
			case app.ParentID:
				mapParentAppSecret[setting.Name] = setting
			case app.ProjectID:
				mapProjectSecret[setting.Name] = setting
			}
		}
	}

	// App inherits ENV vars from the parent app and the project
	appEnv := mapProjectEnv
	for k, v := range mapParentAppEnv {
		appEnv[k] = v
	}
	for k, v := range mapAppEnv {
		appEnv[k] = v
	}

	// TODO: parse all secrets within the ENV values of the app

	// Construct ENV vars into a slice of `k=v`
	for k, v := range appEnv {
		res = append(res, fmt.Sprintf("%s=%s", k, v))
	}
	return res, nil
}
