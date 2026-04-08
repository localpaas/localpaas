package envvarservice

import "github.com/localpaas/localpaas/localpaas_app/entity"

type EnvVar struct {
	*entity.EnvVar
	Errors []string
}

func (env *EnvVar) ToString(sep string) string {
	return env.Key + sep + env.Value
}
