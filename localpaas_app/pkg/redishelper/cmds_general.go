package redishelper

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
)

func Keys(
	ctx context.Context,
	cmder Cmdable,
	pattern string,
) ([]string, error) {
	slice, err := cmder.Keys(ctx, pattern).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, apperrors.Wrap(err)
	}
	return slice, nil
}

func Get[T any](
	ctx context.Context,
	cmder Cmdable,
	key string,
	valueCreator ValueCreator[T],
) (value T, err error) {
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

func MGet[T any](
	ctx context.Context,
	cmder Cmdable,
	keys []string,
	valueCreator ValueCreator[T],
) (values []T, err error) {
	if len(keys) == 0 {
		return values, nil
	}
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
		itemBytes := ParseBytes(item)
		if len(itemBytes) == 0 {
			values = append(values, model.GetData())
			continue
		}
		err = model.RedisUnmarshal(ParseBytes(item))
		if err != nil {
			return nil, apperrors.New(err).WithMsgLog("failed to unmarshal value")
		}
		values = append(values, model.GetData())
	}

	return values, nil
}

func Set[T any](
	ctx context.Context,
	cmder Cmdable,
	key string,
	value Value[T],
	expiration time.Duration,
) (err error) {
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

func SetXX[T any](
	ctx context.Context,
	cmder Cmdable,
	key string,
	value Value[T],
	expiration time.Duration,
) (err error) {
	data, err := value.RedisMarshal()
	if err != nil {
		return apperrors.New(err).WithMsgLog("failed to marshal value")
	}

	_, err = cmder.SetXX(ctx, key, data, expiration).Result()
	if err != nil {
		return apperrors.New(err).WithMsgLog("failed to set value in redis")
	}
	return nil
}

func SetNX[T any](
	ctx context.Context,
	cmder Cmdable,
	key string,
	value Value[T],
	expiration time.Duration,
) (err error) {
	data, err := value.RedisMarshal()
	if err != nil {
		return apperrors.New(err).WithMsgLog("failed to marshal value")
	}

	_, err = cmder.SetNX(ctx, key, data, expiration).Result()
	if err != nil {
		return apperrors.New(err).WithMsgLog("failed to set value in redis")
	}
	return nil
}

func MSet[T any](
	ctx context.Context,
	cmder Cmdable,
	keys []string,
	values []Value[T],
	expiration time.Duration,
) (err error) {
	if len(keys) == 0 {
		return nil
	}
	_, err = cmder.TxPipelined(ctx, func(p redis.Pipeliner) error {
		for i, key := range keys {
			val, err := values[i].RedisMarshal()
			if err != nil {
				return apperrors.New(err).WithMsgLog("failed to marshal value")
			}
			valStr := reflectutil.UnsafeBytesToStr(val)
			if expiration == redis.KeepTTL {
				_, err = p.SetXX(ctx, key, valStr, expiration).Result()
			} else {
				_, err = p.Set(ctx, key, valStr, expiration).Result()
			}
			if err != nil {
				return apperrors.New(err).WithMsgLog("failed to set value in redis")
			}
		}
		return nil
	})
	if err != nil {
		return apperrors.New(err)
	}
	return nil
}

func Del(
	ctx context.Context,
	cmder Cmdable,
	keys ...string,
) (err error) {
	_, err = cmder.Del(ctx, keys...).Result()
	if err != nil {
		return apperrors.New(err).WithMsgLog("failed to delete value in redis")
	}
	return nil
}

func Expire(
	ctx context.Context,
	cmder Cmdable,
	key string,
	expiration time.Duration,
) (err error) {
	_, err = cmder.Expire(ctx, key, expiration).Result()
	if err != nil {
		return apperrors.New(err).WithMsgLog("failed to expire key")
	}
	return nil
}

func ExpireXX(
	ctx context.Context,
	cmder Cmdable,
	key string,
	expiration time.Duration,
) (err error) {
	_, err = cmder.ExpireXX(ctx, key, expiration).Result()
	if err != nil {
		return apperrors.New(err).WithMsgLog("failed to expire key")
	}
	return nil
}
