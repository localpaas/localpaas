package nginxservice

import (
	"context"

	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

const (
	nginxServiceName = "localpaas_nginx"
)

func (s *nginxService) GetNginxSwarmService(ctx context.Context) (*swarm.Service, error) {
	service, err := s.dockerManager.ServiceGetByName(ctx, nginxServiceName)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return service, nil
}
