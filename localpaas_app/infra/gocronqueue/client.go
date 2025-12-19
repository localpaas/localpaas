package gocronqueue

import (
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
)

type Client struct {
	redisClient redis.UniversalClient
	logger      logging.Logger
}

func NewClient(redisClient redis.UniversalClient, logger logging.Logger) (*Client, error) {
	return &Client{
		redisClient: redisClient,
		logger:      logger,
	}, nil
}

func (c *Client) Close() error {
	// Use shared client, so we don't close it
	return nil
}

func (c *Client) ScheduleTask(task *entity.Task, runAt time.Time) error {
	// TODO
	return nil
}

func (c *Client) UnscheduleTask(task *entity.Task) error {
	// TODO
	return nil
}
