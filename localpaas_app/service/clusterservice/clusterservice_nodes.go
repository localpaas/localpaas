package clusterservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

func (s *clusterService) IsMultiNode(ctx context.Context) (bool, error) {
	info, err := s.dockerManager.SystemInfo(ctx)
	if err != nil {
		return false, apperrors.Wrap(err)
	}
	return info.Swarm.Nodes > 1, nil
}
