package entity

import (
	"github.com/localpaas/localpaas/infrastructure/docker"
	"github.com/localpaas/localpaas/localpaas_app/base"
)

func (s *Setting) ParseAppServiceSpec() (*docker.ServiceSpec, error) {
	if s != nil && s.Data != "" && s.Type == base.SettingTypeServiceSpec {
		res := &docker.ServiceSpec{}
		return res, s.parseData(res)
	}
	return nil, nil
}
