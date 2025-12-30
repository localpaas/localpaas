package cacherepository

import (
	"context"
	"fmt"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
	"github.com/localpaas/localpaas/localpaas_app/pkg/redishelper"
)

type UserTokenRepo interface {
	Exist(ctx context.Context, userID, uid string) error
	Set(ctx context.Context, userID, uid string, exp time.Duration) error
	Del(ctx context.Context, userID, uid string) error
	DelAll(ctx context.Context, userID string) error
}

type userTokenRepo struct {
	client rediscache.Client
}

func NewUserTokenRepo(client rediscache.Client) UserTokenRepo {
	return &userTokenRepo{client: client}
}

func (repo *userTokenRepo) Exist(ctx context.Context, userID, uid string) error {
	key := repo.formatKey(userID, uid)
	count, err := repo.client.Exists(ctx, key).Result()
	if err != nil {
		return apperrors.New(err)
	}
	if count == 0 {
		return apperrors.NewNotFoundNT(key)
	}
	return nil
}

func (repo *userTokenRepo) Set(ctx context.Context, userID, uid string, exp time.Duration) error {
	err := redishelper.Set(ctx, repo.client, repo.formatKey(userID, uid),
		redishelper.NewJSONValue(""), exp)
	if err != nil {
		return apperrors.New(err)
	}
	return nil
}

func (repo *userTokenRepo) Del(ctx context.Context, userID, uid string) error {
	err := redishelper.Del(ctx, repo.client, repo.formatKey(userID, uid))
	if err != nil {
		return apperrors.New(err)
	}
	return nil
}

func (repo *userTokenRepo) DelAll(ctx context.Context, userID string) error {
	keys, err := repo.client.Keys(ctx, repo.formatKey(userID, "*")).Result()
	if err != nil {
		return apperrors.New(err)
	}
	if len(keys) == 0 {
		return nil
	}
	err = redishelper.Del(ctx, repo.client, keys...)
	if err != nil {
		return apperrors.New(err)
	}
	return nil
}

func (repo *userTokenRepo) formatKey(userID, uid string) string {
	return fmt.Sprintf("user:%s:token:%s", userID, uid)
}
