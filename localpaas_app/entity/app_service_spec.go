package entity

import (
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/services/docker"
)

func (s *Setting) ParseAppServiceSpec() (*docker.ServiceSpec, error) {
	res := &docker.ServiceSpec{}
	if s != nil && s.Data != "" && s.Type == base.SettingTypeServiceSpec {
		return res, s.parseData(res)
	}
	return res, nil
}
