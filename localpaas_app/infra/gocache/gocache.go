package gocache

import (
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/pkg/timeutil"
)

const (
	defaultExpiration = time.Minute * 5
	cleanupInterval   = 4 * time.Hour
)

const (
	NoExpiration time.Duration = cache.NoExpiration
)

type Cache struct {
	client *cache.Cache
}

func NewCache() *Cache {
	return &Cache{
		client: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (c *Cache) Get(key string) (any, error) {
	val, exists := c.client.Get(key)
	if !exists {
		return nil, apperrors.NewNotFoundNT(key)
	}
	return val, nil
}

func (c *Cache) Set(key string, value any, expiration time.Duration) error {
	c.client.Set(key, value, expiration)
	return nil
}

func (c *Cache) Del(key string) error {
	c.client.Delete(key)
	return nil
}

func (c *Cache) TTL(key string) (time.Duration, error) {
	_, t, exists := c.client.GetWithExpiration(key)
	if !exists {
		return 0, nil
	}
	if t.IsZero() {
		return -1, nil
	}
	return t.Sub(timeutil.NowUTC()), nil
}
