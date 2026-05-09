package healthcheckserviceimpl

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/service/healthcheckservice"
)

type service struct {
	logger logging.Logger
}

func New(
	logger logging.Logger,
) healthcheckservice.Service {
	return &service{
		logger: logger,
	}
}
