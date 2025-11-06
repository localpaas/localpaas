package entity

import (
	"github.com/localpaas/localpaas/localpaas_app/base"
)

type ProjectSettings struct {
	Test string `json:"test"`
}

func (s *Setting) ParseProjectSettings() (*ProjectSettings, error) {
	if s != nil && s.Type == base.SettingTypeProject {
		res := &ProjectSettings{}
		return res, s.parseData(res)
	}
	return nil, nil
}
