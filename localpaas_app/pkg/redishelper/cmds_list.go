package redishelper

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
)

func RPush[T any](
	ctx context.Context,
	cmder Cmdable,
	key string,
	values ...T,
) error {
	if len(values) == 0 {
		return nil
	}
	data := make([]any, 0, len(values))
	for _, v := range values {
		item, err := jsonMarshal(v)
		if err != nil {
			return apperrors.New(err).WithMsgLog("failed to marshal value")
		}
		data = append(data, item)
	}
	_, err := cmder.RPush(ctx, key, data...).Result()
	if err != nil {
		return apperrors.New(err).WithMsgLog("failed to push values to a list")
	}
	return nil
}

func LRange[T any](
	ctx context.Context,
	cmder Cmdable,
	key string,
	start, stop int64,
) (values []T, err error) {
	data, err := cmder.LRange(ctx, key, start, stop).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, apperrors.Wrap(err)
	}

	return unmarshalStrSlice[T](data...)
}

func BLPop[T any](
	ctx context.Context,
	cmder Cmdable,
	keys []string,
	timeout time.Duration,
) (values map[string]T, err error) {
	strSlice, err := cmder.BLPop(ctx, timeout, keys...).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, apperrors.Wrap(err)
	}

	values = make(map[string]T)
	i := 0
	for i < len(strSlice) {
		key := strSlice[i]
		i++
		valStr := strSlice[i]
		i++

		var val T
		err := jsonUnmarshal(reflectutil.UnsafeStrToBytes(valStr), &val)
		if err != nil {
			return nil, apperrors.New(err).WithMsgLog("failed to unmarshal value")
		}
		values[key] = val
	}

	return values, nil
}
