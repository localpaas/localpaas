package entity

import (
	"github.com/localpaas/localpaas/localpaas_app/base"
)

type AppDeploymentSettings struct {
	Test string `json:"test"`

	// NOTE: for storing current containing setting only
	Setting *Setting `json:"-"`
}

func (s *Setting) ParseAppDeploymentSettings() (*AppDeploymentSettings, error) {
	res := &AppDeploymentSettings{Setting: s}
	if s != nil && s.Data != "" && s.Type == base.SettingTypeDeployment {
		return res, s.parseData(res)
	}
	return nil, nil
}
