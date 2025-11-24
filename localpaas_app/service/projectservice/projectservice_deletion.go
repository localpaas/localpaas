package projectservice

import (
	"context"
	"sync"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

func (s *projectService) DeleteProject(ctx context.Context, project *entity.Project) error {
	// Remove all apps
	var wg sync.WaitGroup
	for _, app := range project.Apps {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.appService.DeleteApp(ctx, app)
			// NOTE: it's hard to rollback, maybe we only show the errors if there is any
		}()
	}
	wg.Wait()

	// Remove project network
	err := s.networkService.RemoveProjectNetwork(ctx, project)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
