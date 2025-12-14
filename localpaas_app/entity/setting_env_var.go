package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentEnvVarsVersion = 1
)

type EnvVars struct {
	Data []*EnvVar `json:"data"`
}

type EnvVar struct {
	Key        string `json:"k"`
	Value      string `json:"v"`
	IsBuildEnv bool   `json:"isBuildEnv,omitempty"`
}

func (s *Setting) AsEnvVars() (*EnvVars, error) {
	return parseSettingAs(s, base.SettingTypeEnvVar, func() *EnvVars { return &EnvVars{} })
}

func (s *Setting) MustAsEnvVars() *EnvVars {
	return gofn.Must(s.AsEnvVars())
}
