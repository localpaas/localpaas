package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/services/docker"
)

func (s *Setting) AsAppServiceSpec() (*docker.ServiceSpec, error) {
	if s.parsedData != nil {
		res, ok := s.parsedData.(*docker.ServiceSpec)
		if !ok {
			return nil, apperrors.NewTypeInvalid()
		}
		return res, nil
	}
	res := &docker.ServiceSpec{}
	if s.Data != "" && s.Type == base.SettingTypeServiceSpec {
		if err := s.parseData(res); err != nil {
			return nil, apperrors.Wrap(err)
		}
	}
	return res, nil
}

func (s *Setting) MustAsAppServiceSpec() *docker.ServiceSpec {
	return gofn.Must(s.AsAppServiceSpec())
}
