package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
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
	if s.parsedData != nil {
		res, ok := s.parsedData.(*EnvVars)
		if !ok {
			return nil, apperrors.NewTypeInvalid()
		}
		return res, nil
	}
	res := &EnvVars{}
	if s.Data != "" && s.Type == base.SettingTypeEnvVar {
		if err := s.parseData(res); err != nil {
			return nil, apperrors.Wrap(err)
		}
	}
	return res, nil
}

func (s *Setting) MustAsEnvVars() *EnvVars {
	return gofn.Must(s.AsEnvVars())
}
