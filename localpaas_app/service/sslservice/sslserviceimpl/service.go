package sslserviceimpl

import "github.com/localpaas/localpaas/localpaas_app/service/sslservice"

func New() sslservice.Service {
	return &service{}
}

type service struct {
}
