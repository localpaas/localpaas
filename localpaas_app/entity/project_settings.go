package entity

import (
	"github.com/localpaas/localpaas/localpaas_app/base"
)

type ProjectSettings struct {
	Test string `json:"test"`

	// NOTE: for storing current containing setting only
	Setting *Setting `json:"-"`
}

func (s *Setting) ParseProjectSettings() (*ProjectSettings, error) {
	res := &ProjectSettings{Setting: s}
	if s != nil && s.Data != "" && s.Type == base.SettingTypeProject {
		return res, s.parseData(res)
	}
	return res, nil
}
