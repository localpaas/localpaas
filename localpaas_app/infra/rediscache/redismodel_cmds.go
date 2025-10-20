package rediscache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/pkg/reflectutil"
)

func Get[T any](ctx context.Context, cmder Cmdable, key string, valueCreator ValueCreator[T]) (value T, err error) {
	valueStr, err := cmder.Get(ctx, key).Result()
	if err != nil && errors.Is(err, redis.Nil) {
		return value, apperrors.NewNotFoundNT(key)
	}

	model := valueCreator(value)
	err = model.RedisUnmarshal(reflectutil.UnsafeStrToBytes(valueStr))
	if err != nil {
		return value, apperrors.New(err).WithMsgLog("failed to unmarshal value")
	}
	return model.GetData(), nil
}

func Set[T any](ctx context.Context, cmder Cmdable, key string, value Value[T], expiration time.Duration) (err error) {
	data, err := value.RedisMarshal()
	if err != nil {
		return apperrors.New(err).WithMsgLog("failed to marshal value")
	}

	_, err = cmder.Set(ctx, key, data, expiration).Result()
	if err != nil {
		return apperrors.New(err).WithMsgLog("failed to set value in redis")
	}
	return nil
}

func Del(ctx context.Context, cmder Cmdable, key string) (err error) {
	_, err = cmder.Del(ctx, key).Result()
	if err != nil {
		return apperrors.New(err).WithMsgLog("failed to delete value in redis")
	}
	return nil
}
