package taskqueue

import (
	"sync"

	"github.com/hibiken/asynq"

	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
)

var (
	clientOnce sync.Once
	client     *asynq.Client
)

func StartClient(redisClient rediscache.Client) error {
	clientOnce.Do(func() {
		client = asynq.NewClient(redisConnOpt{client: redisClient})
	})
	return nil
}

func StopClient() error {
	if client != nil {
		return client.Close() //nolint:wrapcheck
	}
	return nil
}
