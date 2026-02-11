package taskservice

import (
	"context"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

func (s *taskService) CreateLock(
	ctx context.Context,
	key string,
	exp time.Duration,
) (success bool, releaser func(), err error) {
	success, err = s.redisClient.SetNX(ctx, key, "1", exp).Result()
	if err != nil {
		return false, nil, apperrors.Wrap(err)
	}
	return success, func() {
		_, _ = s.redisClient.Del(ctx, key).Result()
	}, nil
}
