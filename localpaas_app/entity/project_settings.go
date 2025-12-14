package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentProjectSettingsVersion = 1
)

type ProjectSettings struct {
	Test string `json:"test"`
}

func (s *Setting) AsProjectSettings() (*ProjectSettings, error) {
	if s.parsedData != nil {
		res, ok := s.parsedData.(*ProjectSettings)
		if !ok {
			return nil, apperrors.NewTypeInvalid()
		}
		return res, nil
	}
	res := &ProjectSettings{}
	if s.Data != "" && s.Type == base.SettingTypeProject {
		if err := s.parseData(res); err != nil {
			return nil, apperrors.Wrap(err)
		}
	}
	return res, nil
}

func (s *Setting) MustAsProjectSettings() *ProjectSettings {
	return gofn.Must(s.AsProjectSettings())
}
