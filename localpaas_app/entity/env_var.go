package entity

import "github.com/localpaas/localpaas/localpaas_app/base"

type EnvVars struct {
	Data []*EnvVar `json:"data"`
}

type EnvVar struct {
	Key        string `json:"k"`
	Value      string `json:"v"`
	IsBuildEnv bool   `json:"isBuildEnv,omitempty"`
}

func (s *Setting) ParseEnvVars() (*EnvVars, error) {
	if s != nil && s.Type == base.SettingTypeEnvVar {
		res := &EnvVars{}
		return res, s.parseData(res)
	}
	return nil, nil
}
