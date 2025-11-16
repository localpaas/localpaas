package rediscache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
)

func Keys(ctx context.Context, cmder Cmdable, pattern string) ([]string, error) {
	slice, err := cmder.Keys(ctx, pattern).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, apperrors.Wrap(err)
	}
	return slice, nil
}

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

func MGet[T any](ctx context.Context, cmder Cmdable, keys []string, valueCreator ValueCreator[T]) (
	values []T, err error) {
	slice, err := cmder.MGet(ctx, keys...).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, apperrors.Wrap(err)
	}

	var valDefault T
	for _, item := range slice {
		model := valueCreator(valDefault)
		err = model.RedisUnmarshal(ParseBytes(item))
		if err != nil {
			return nil, apperrors.New(err).WithMsgLog("failed to unmarshal value")
		}
		values = append(values, model.GetData())
	}

	return values, nil
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

func MSet[T any](ctx context.Context, cmder Cmdable, keys []string, values []Value[T],
	expiration time.Duration) (err error) {
	setValues := make([]any, 0, len(values)*2) //nolint:mnd
	for i, v := range values {
		item, err := v.RedisMarshal()
		if err != nil {
			return apperrors.New(err).WithMsgLog("failed to marshal value")
		}
		setValues = append(setValues, keys[i], reflectutil.UnsafeBytesToStr(item))
	}

	_, err = cmder.TxPipelined(ctx, func(p redis.Pipeliner) error {
		_, err = p.MSet(ctx, setValues...).Result()
		if err != nil {
			return apperrors.New(err).WithMsgLog("failed to set values in redis")
		}
		if expiration == 0 {
			return nil
		}
		for _, key := range keys {
			_, err = p.Expire(ctx, key, expiration).Result()
			if err != nil {
				return apperrors.New(err).WithMsgLog("failed to set expiration in redis")
			}
		}
		return nil
	})
	if err != nil {
		return apperrors.New(err)
	}
	return nil
}

func Del(ctx context.Context, cmder Cmdable, keys ...string) (err error) {
	_, err = cmder.Del(ctx, keys...).Result()
	if err != nil {
		return apperrors.New(err).WithMsgLog("failed to delete value in redis")
	}
	return nil
}
