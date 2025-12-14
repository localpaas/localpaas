package gocronqueue

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
)

type Client struct {
	redisClient rediscache.Client
	logger      logging.Logger
}

func NewClient(redisClient rediscache.Client, logger logging.Logger) (*Client, error) {
	return &Client{
		redisClient: redisClient,
		logger:      logger,
	}, nil
}

func (c *Client) Close() error {
	return nil
}
