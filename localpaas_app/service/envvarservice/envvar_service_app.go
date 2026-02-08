package envvarservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
)

type EnvVar struct {
	*entity.EnvVar
	Errors []string
}

func (env *EnvVar) ToString(sep string) string {
	return env.Key + sep + env.Value
}

func (s *envVarService) BuildAppEnvVars(
	ctx context.Context,
	db database.IDB,
	app *entity.App,
	buildPhase bool,
) (res []*EnvVar, err error) {
	vars, secrets, err := s.LoadAppEnvVarsAndSecrets(ctx, db, app, true, true, buildPhase)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	// App inherits ENV vars in order from the parent app, then the project, and then from global
	envStore := vars.FinalEnvVars()
	secretStore := secrets.FinalSecrets()

	// Construct result
	res = make([]*EnvVar, 0, len(envStore))
	for _, env := range envStore {
		res = append(res, &EnvVar{EnvVar: env})
	}

	// Process all references within the ENV values
	for _, env := range res {
		s.processRefs(env, envStore, secretStore)
	}

	return res, nil
}

func (s *envVarService) ProcessEnvRefs(
	ctx context.Context,
	db database.IDB,
	app *entity.App,
	envVars []*entity.EnvVar,
	loadEnvVars bool,
	loadSecrets bool,
	buildPhase bool,
) (res []*EnvVar, err error) {
	if len(envVars) == 0 {
		return nil, nil
	}
	vars, secrets, err := s.LoadAppEnvVarsAndSecrets(ctx, db, app, loadEnvVars, loadSecrets, buildPhase)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	// Construct result
	res = make([]*EnvVar, 0, len(envVars))
	for _, env := range envVars {
		res = append(res, &EnvVar{EnvVar: env})
	}

	// App inherits ENV vars in order from the parent app, then the project, and then from global
	envStore := vars.FinalEnvVars()
	// Update the envStore with the input values
	for _, env := range envVars {
		envStore[env.Key] = env
	}
	secretStore := secrets.FinalSecrets()
	// Process all references within the ENV values
	for _, env := range res {
		s.processRefs(env, envStore, secretStore)
	}
	return res, nil
}

type AppEnvVars struct {
	App       map[string]*entity.EnvVar
	ParentApp map[string]*entity.EnvVar
	Project   map[string]*entity.EnvVar
	Global    map[string]*entity.EnvVar
}

func (ev *AppEnvVars) FinalEnvVars() map[string]*entity.EnvVar {
	res := make(map[string]*entity.EnvVar, 20) //nolint
	for k, v := range ev.Global {
		res[k] = v
	}
	for k, v := range ev.Project {
		res[k] = v
	}
	for k, v := range ev.ParentApp {
		res[k] = v
	}
	for k, v := range ev.App {
		res[k] = v
	}
	return res
}

type AppSecrets struct {
	App       map[string]*entity.Secret
	ParentApp map[string]*entity.Secret
	Project   map[string]*entity.Secret
	Global    map[string]*entity.Secret
}

func (ev *AppSecrets) FinalSecrets() map[string]*entity.Secret {
	res := make(map[string]*entity.Secret, 20) //nolint
	for k, v := range ev.Global {
		res[k] = v
	}
	for k, v := range ev.Project {
		res[k] = v
	}
	for k, v := range ev.ParentApp {
		res[k] = v
	}
	for k, v := range ev.App {
		res[k] = v
	}
	return res
}

//nolint:gocognit
func (s *envVarService) LoadAppEnvVarsAndSecrets(
	ctx context.Context,
	db database.IDB,
	app *entity.App,
	loadEnvVars bool,
	loadSecrets bool,
	buildPhase bool,
) (envVars *AppEnvVars, secrets *AppSecrets, err error) {
	var settingTypes []base.SettingType
	if loadEnvVars {
		settingTypes = append(settingTypes, base.SettingTypeEnvVar)
		//nolint:mnd
		envVars = &AppEnvVars{
			App:       make(map[string]*entity.EnvVar, 20),
			ParentApp: make(map[string]*entity.EnvVar, 20),
			Project:   make(map[string]*entity.EnvVar, 20),
			Global:    make(map[string]*entity.EnvVar, 20),
		}
	}
	if loadSecrets {
		settingTypes = append(settingTypes, base.SettingTypeSecret)
		//nolint:mnd
		secrets = &AppSecrets{
			App:       make(map[string]*entity.Secret, 10),
			ParentApp: make(map[string]*entity.Secret, 10),
			Project:   make(map[string]*entity.Secret, 10),
			Global:    make(map[string]*entity.Secret, 10),
		}
	}
	if len(settingTypes) == 0 {
		return nil, nil, nil
	}

	settings, _, err := s.settingRepo.ListByAppObject(ctx, db, app, nil,
		bunex.SelectWhereIn("setting.type IN (?)", settingTypes...),
		bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
	)
	if err != nil {
		return nil, nil, apperrors.Wrap(err)
	}

	if len(settings) == 0 {
		return envVars, secrets, nil
	}

	for _, setting := range settings {
		if setting.Type == base.SettingTypeEnvVar {
			vars := setting.MustAsEnvVars()
			switch {
			case setting.ObjectID == app.ID:
				for _, env := range vars.Data {
					if env.IsBuildEnv == buildPhase {
						envVars.App[env.Key] = env
					}
				}
			case app.ParentID != "" && setting.ObjectID == app.ParentID:
				for _, env := range vars.Data {
					if env.IsBuildEnv == buildPhase {
						envVars.ParentApp[env.Key] = env
					}
				}
			case setting.ObjectID == "":
				for _, env := range vars.Data {
					if env.IsBuildEnv == buildPhase {
						envVars.Global[env.Key] = env
					}
				}
			default:
				for _, env := range vars.Data {
					if env.IsBuildEnv == buildPhase {
						envVars.Project[env.Key] = env
					}
				}
			}
			continue
		}

		if setting.Type == base.SettingTypeSecret {
			secret := setting.MustAsSecret() // decryption takes time, so do it when needed only
			switch {
			case setting.ObjectID == app.ID:
				secrets.App[setting.Name] = secret
			case app.ParentID != "" && setting.ObjectID == app.ParentID:
				secrets.ParentApp[setting.Name] = secret
			case setting.ObjectID == "":
				secrets.Global[setting.Name] = secret
			default:
				secrets.Project[setting.Name] = secret
			}
		}
	}

	return envVars, secrets, nil
}
