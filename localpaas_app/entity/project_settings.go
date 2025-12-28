package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentProjectSettingsVersion = 1
)

type ProjectSettings struct {
	Test string `json:"test"`
}

func (s *Setting) AsProjectSettings() (*ProjectSettings, error) {
	return parseSettingAs(s, base.SettingTypeProject, func() *ProjectSettings { return &ProjectSettings{} })
}

func (s *Setting) MustAsProjectSettings() *ProjectSettings {
	return gofn.Must(s.AsProjectSettings())
}
