package cacherepository

import (
	"context"
	"fmt"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity/cacheentity"
	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
	"github.com/localpaas/localpaas/localpaas_app/pkg/redishelper"
)

type TaskControlRepo interface {
	Push(ctx context.Context, taskID string, taskControl *cacheentity.TaskControl) error
}

type taskControlRepo struct {
	client rediscache.Client
}

func NewTaskControlRepo(client rediscache.Client) TaskControlRepo {
	return &taskControlRepo{client: client}
}

func (repo *taskControlRepo) Push(
	ctx context.Context,
	taskID string,
	taskControl *cacheentity.TaskControl,
) error {
	key := repo.formatKey(taskID)
	err := redishelper.RPush(ctx, repo.client, key, taskControl)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (repo *taskControlRepo) formatKey(taskID string) string {
	return fmt.Sprintf("task:%s:ctrl", taskID)
}
