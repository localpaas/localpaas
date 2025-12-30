package realtimelog

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/pkg/redishelper"
)

const (
	storeExpiration = 4 * time.Hour
)

type Store struct {
	redisClient redis.UniversalClient
	key         string
	initialized bool
	storeLocal  bool
	mu          sync.RWMutex
	frames      []*LogFrame
}

func (s *Store) Add(ctx context.Context, frames ...*LogFrame) error {
	if s.storeLocal {
		s.mu.Lock()
		s.frames = append(s.frames, frames...)
		s.mu.Unlock()
	}

	// Store log data in redis
	err := redishelper.RPush(ctx, s.redisClient, s.key, redishelper.NewJSONValues(frames)...)
	if err != nil {
		return apperrors.New(err).WithMsgLog("failed to push log frames to redis")
	}

	if !s.initialized {
		s.redisClient.ExpireXX(ctx, s.key, storeExpiration)
		s.initialized = true
	}

	// Notify consumers about the new data
	_, err = s.redisClient.Publish(ctx, s.key, buildMessage(CommandNewData)).Result()
	if err != nil {
		return apperrors.New(err).WithMsgLog("failed to notify consumers about the new data")
	}

	return nil
}

func (s *Store) GetData(ctx context.Context, fromIndex int64) ([]*LogFrame, error) {
	if s.storeLocal {
		return s.GetLocalData(ctx, fromIndex)
	}
	return s.GetRemoteData(ctx, fromIndex)
}

func (s *Store) GetLocalData(ctx context.Context, fromIndex int64) ([]*LogFrame, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if fromIndex >= int64(len(s.frames)) {
		return nil, nil
	}
	return s.frames[fromIndex:], nil
}

func (s *Store) GetRemoteData(ctx context.Context, fromIndex int64) ([]*LogFrame, error) {
	frames, err := redishelper.LRange(ctx, s.redisClient, s.key, fromIndex, -1,
		redishelper.JSONValueCreator[*LogFrame])
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return frames, nil
}

func (s *Store) Close() (err error) {
	ctx := context.Background()
	defer func() {
		removeLocalStore(s.key)
	}()

	// Send close-msg to consumers
	_, e := s.redisClient.Publish(ctx, s.key, buildMessage(CommandClosed)).Result()
	if e != nil {
		err = errors.Join(err, apperrors.New(err).WithMsgLog("failed to notify consumers"))
	}
	// Delete log data in redis
	e = redishelper.Del(ctx, s.redisClient, s.key)
	if e != nil {
		err = errors.Join(err, apperrors.New(err).WithMsgLog("failed to remove data from redis"))
	}

	return err
}

func NewStore(
	key string,
	storeLocal bool,
	redisClient redis.UniversalClient,
) *Store {
	s := &Store{
		redisClient: redisClient,
		key:         key,
		storeLocal:  storeLocal,
	}
	if storeLocal {
		addLocalStore(key, s)
		s.frames = make([]*LogFrame, 0, 100) //nolint:mnd
	}
	return s
}
