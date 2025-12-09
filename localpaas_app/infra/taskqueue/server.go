package taskqueue

import (
	"sync"

	"github.com/hibiken/asynq"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
)

var (
	serverOnce sync.Once
	server     *asynq.Server
)

type redisConnOpt struct {
	client any
}

func (r redisConnOpt) MakeRedisClient() any {
	return r.client
}

type TaskHandlerMux func(mux *asynq.ServeMux) error

//nolint:mnd
func StartServer(client rediscache.Client, taskHandler TaskHandlerMux, logger logging.Logger) error {
	serverOnce.Do(func() {
		server = asynq.NewServer(
			redisConnOpt{client: client},
			asynq.Config{
				// Specify how many concurrent workers to use
				Concurrency: 10,
				// Optionally specify multiple queues with different priority
				Queues: map[string]int{
					"critical": 5,
					"default":  3,
					"low":      2,
				},
			},
		)
	})

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	err := taskHandler(mux)
	if err != nil {
		return apperrors.Wrap(err)
	}

	go func() {
		if err := server.Run(mux); err != nil {
			logger.Errorf("failed to start task queue server: %v", err)
		}
	}()

	return nil
}

func StopServer() error {
	if server != nil {
		server.Stop()
		server.Shutdown()
	}
	return nil
}
