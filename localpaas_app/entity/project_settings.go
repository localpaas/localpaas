package entity

import (
	"github.com/localpaas/localpaas/localpaas_app/base"
)

type ProjectSettings struct {
	Test string `json:"test"`
}

func (s *Setting) ParseProjectSettings() (*ProjectSettings, error) {
	res := &ProjectSettings{}
	if s != nil && s.Type == base.SettingTypeProject {
		return res, s.parseData(res)
	}
	return res, nil
}
