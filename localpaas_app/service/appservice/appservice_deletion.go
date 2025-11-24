package appservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

func (s *appService) DeleteApp(ctx context.Context, app *entity.App) error {
	// Remove service for the app in docker swarm
	err := s.dockerManager.ServiceRemove(ctx, app.ServiceID)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Remove app config from nginx
	err = s.nginxService.RemoveAppConfig(ctx, app)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
