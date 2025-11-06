package entity

import (
	"github.com/localpaas/localpaas/localpaas_app/base"
)

type AppDeploymentSettings struct {
	Test string `json:"test"`
}

func (s *Setting) ParseAppDeploymentSettings() (*AppDeploymentSettings, error) {
	if s != nil && s.Type == base.SettingTypeDeployment {
		res := &AppDeploymentSettings{}
		return res, s.parseData(res)
	}
	return nil, nil
}
