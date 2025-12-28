package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/services/docker"
)

func (s *Setting) AsAppServiceSpec() (*docker.ServiceSpec, error) {
	return parseSettingAs(s, base.SettingTypeServiceSpec, func() *docker.ServiceSpec { return &docker.ServiceSpec{} })
}

func (s *Setting) MustAsAppServiceSpec() *docker.ServiceSpec {
	return gofn.Must(s.AsAppServiceSpec())
}
