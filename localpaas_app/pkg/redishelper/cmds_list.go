package redishelper

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
)

func RPush[T any](
	ctx context.Context,
	cmder Cmdable,
	key string,
	values ...Value[T],
) error {
	if len(values) == 0 {
		return nil
	}
	data := make([]any, 0, len(values))
	for _, v := range values {
		item, err := v.RedisMarshal()
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
	valueCreator ValueCreator[T],
) (values []T, err error) {
	strSlice, err := cmder.LRange(ctx, key, start, stop).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, apperrors.Wrap(err)
	}

	var valDefault T
	for _, item := range strSlice {
		model := valueCreator(valDefault)
		if len(item) == 0 {
			values = append(values, model.GetData())
			continue
		}
		err = model.RedisUnmarshal(reflectutil.UnsafeStrToBytes(item))
		if err != nil {
			return nil, apperrors.New(err).WithMsgLog("failed to unmarshal value")
		}
		values = append(values, model.GetData())
	}

	return values, nil
}
