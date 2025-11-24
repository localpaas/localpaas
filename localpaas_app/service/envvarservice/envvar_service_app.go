package envvarservice

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
)

type EnvVar struct {
	Key   string
	Value string
	Error string
}

func (env *EnvVar) ToString(sep string) string {
	return env.Key + sep + env.Value
}

//nolint:gocognit
func (s *envVarService) BuildAppEnv(ctx context.Context, db database.IDB, app *entity.App, buildPhase bool) (
	res []*EnvVar, err error) {
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

	mapAppEnv := make(map[string]string, 20)                  //nolint
	mapParentAppEnv := make(map[string]string, 20)            //nolint
	mapProjectEnv := make(map[string]string, 20)              //nolint
	mapAppSecret := make(map[string]*entity.Secret, 10)       //nolint
	mapParentAppSecret := make(map[string]*entity.Secret, 10) //nolint
	mapProjectSecret := make(map[string]*entity.Secret, 10)   //nolint

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
			secret, err := setting.ParseSecret(false) // decryption takes time, so do it when needed only
			if err != nil {
				return nil, apperrors.Wrap(err)
			}
			switch setting.ObjectID {
			case app.ID:
				mapAppSecret[setting.Name] = secret
			case app.ParentID:
				mapParentAppSecret[setting.Name] = secret
			case app.ProjectID:
				mapProjectSecret[setting.Name] = secret
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

	// Construct result
	res = make([]*EnvVar, 0, len(appEnv))
	for k, v := range appEnv {
		res = append(res, &EnvVar{Key: k, Value: v})
	}

	// Process all references within the ENV values
	secretStores := []map[string]*entity.Secret{mapAppSecret, mapParentAppSecret, mapProjectSecret}
	envStores := []map[string]string{mapAppEnv, mapParentAppEnv, mapProjectEnv}
	for _, env := range res {
		s.ProcessEnvVarRefs(env, secretStores, envStores)
	}

	return res, nil
}
