package appserviceimpl

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

func (s *service) DeleteApp(ctx context.Context, app *entity.App) error {
	// Remove service for the app in docker swarm
	_, err := s.dockerManager.ServiceRemove(ctx, app.ServiceID)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Remove app config from traefik
	err = s.traefikService.RemoveAppConfig(ctx, app, nil)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
