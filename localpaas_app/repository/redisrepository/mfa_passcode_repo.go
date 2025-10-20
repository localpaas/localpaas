package redisrepository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/localpaas/localpaas/localpaas_app/entity/redisentity"
	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
)

type MFAPasscodeRepo interface {
	Get(ctx context.Context, userID string) (*redisentity.MFAPasscode, error)
	TTL(ctx context.Context, userID string) (time.Duration, error)
	Set(ctx context.Context, userID string, passcode *redisentity.MFAPasscode, exp time.Duration) error
	IncrAttempts(ctx context.Context, userID string, passcode *redisentity.MFAPasscode) error
	Del(ctx context.Context, userID string) error
}

type mfaPasscodeRepo struct {
	client rediscache.Client
}

func NewMFAPasscodeRepo(client rediscache.Client) MFAPasscodeRepo {
	return &mfaPasscodeRepo{client: client}
}

func (repo *mfaPasscodeRepo) Get(ctx context.Context, userID string) (*redisentity.MFAPasscode, error) {
	//nolint:wrapcheck
	return rediscache.Get(ctx, repo.client, repo.formatKey(userID), rediscache.NewJSONValue[*redisentity.MFAPasscode])
}

func (repo *mfaPasscodeRepo) TTL(ctx context.Context, userID string) (time.Duration, error) {
	//nolint:wrapcheck
	return repo.client.TTL(ctx, repo.formatKey(userID)).Result()
}

func (repo *mfaPasscodeRepo) Set(ctx context.Context, userID string, passcode *redisentity.MFAPasscode,
	exp time.Duration) error {
	//nolint:wrapcheck
	return rediscache.Set(ctx, repo.client, repo.formatKey(userID), rediscache.NewJSONValue(passcode), exp)
}

func (repo *mfaPasscodeRepo) IncrAttempts(ctx context.Context, userID string,
	passcode *redisentity.MFAPasscode) error {
	passcode.Attempts++
	// Only set value, not set expiration to keep the current expiration value
	//nolint:wrapcheck
	return rediscache.Set(ctx, repo.client, repo.formatKey(userID), rediscache.NewJSONValue(passcode), redis.KeepTTL)
}

func (repo *mfaPasscodeRepo) Del(ctx context.Context, userID string) error {
	//nolint:wrapcheck
	return rediscache.Del(ctx, repo.client, repo.formatKey(userID))
}

func (repo *mfaPasscodeRepo) formatKey(userID string) string {
	return fmt.Sprintf("mfa-passcode:%s", userID)
}
