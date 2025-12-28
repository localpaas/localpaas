package cacherepository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity/cacheentity"
	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
)

type DeploymentInfoRepo interface {
	Get(ctx context.Context, deploymentID string) (*cacheentity.DeploymentInfo, error)
	MGet(ctx context.Context, deploymentIDs []string) (map[string]*cacheentity.DeploymentInfo, error)
	GetAllOfApp(ctx context.Context, appID string) (map[string]*cacheentity.DeploymentInfo, error)
	GetAll(ctx context.Context) (map[string]*cacheentity.DeploymentInfo, error)
	Set(ctx context.Context, deploymentID string, deploymentInfo *cacheentity.DeploymentInfo, exp time.Duration) error
	Update(ctx context.Context, deploymentID string, deploymentInfo *cacheentity.DeploymentInfo) error
	CancelAllOfApp(ctx context.Context, appID string) error
	Del(ctx context.Context, deploymentID string) error
}

type deploymentInfoRepo struct {
	client rediscache.Client
}

func NewDeploymentInfoRepo(client rediscache.Client) DeploymentInfoRepo {
	return &deploymentInfoRepo{client: client}
}

func (repo *deploymentInfoRepo) Get(
	ctx context.Context,
	deploymentID string,
) (*cacheentity.DeploymentInfo, error) {
	resp, err := rediscache.Get(ctx, repo.client, repo.formatKey(deploymentID),
		rediscache.NewJSONValue[*cacheentity.DeploymentInfo])
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}

func (repo *deploymentInfoRepo) MGet(
	ctx context.Context,
	deploymentIDs []string,
) (map[string]*cacheentity.DeploymentInfo, error) {
	if len(deploymentIDs) == 0 {
		return nil, nil
	}
	keys := make([]string, 0, len(deploymentIDs))
	for _, deploymentID := range deploymentIDs {
		keys = append(keys, repo.formatKey(deploymentID))
	}
	return repo.mGet(ctx, keys)
}

func (repo *deploymentInfoRepo) GetAllOfApp(
	ctx context.Context,
	appID string,
) (map[string]*cacheentity.DeploymentInfo, error) {
	keys, err := rediscache.Keys(ctx, repo.client, repo.formatKey("*"))
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if len(keys) == 0 {
		return nil, nil
	}
	deployments, err := repo.mGet(ctx, keys)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	appDeployments := make(map[string]*cacheentity.DeploymentInfo, len(deployments))
	for k, deployment := range deployments {
		if deployment.AppID != appID {
			continue
		}
		appDeployments[k] = deployment
	}
	return appDeployments, nil
}

func (repo *deploymentInfoRepo) GetAll(
	ctx context.Context,
) (map[string]*cacheentity.DeploymentInfo, error) {
	keys, err := rediscache.Keys(ctx, repo.client, repo.formatKey("*"))
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if len(keys) == 0 {
		return nil, nil
	}
	return repo.mGet(ctx, keys)
}

func (repo *deploymentInfoRepo) mGet(ctx context.Context, keys []string) (
	map[string]*cacheentity.DeploymentInfo, error) {
	resp, err := rediscache.MGet(ctx, repo.client, keys, rediscache.NewJSONValue[*cacheentity.DeploymentInfo])
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	result := make(map[string]*cacheentity.DeploymentInfo, len(resp))
	for _, info := range resp {
		if info != nil {
			result[info.ID] = info
		}
	}
	return result, nil
}

func (repo *deploymentInfoRepo) Set(
	ctx context.Context,
	deploymentID string,
	deploymentInfo *cacheentity.DeploymentInfo,
	exp time.Duration,
) error {
	err := rediscache.Set(ctx, repo.client, repo.formatKey(deploymentID),
		rediscache.NewJSONValue(deploymentInfo), exp)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (repo *deploymentInfoRepo) Update(
	ctx context.Context,
	deploymentID string,
	deploymentInfo *cacheentity.DeploymentInfo,
) error {
	err := rediscache.SetXX(ctx, repo.client, repo.formatKey(deploymentID),
		rediscache.NewJSONValue(deploymentInfo), redis.KeepTTL)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (repo *deploymentInfoRepo) CancelAllOfApp(ctx context.Context, appID string) error {
	deployments, err := repo.GetAllOfApp(ctx, appID)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if len(deployments) == 0 {
		return nil
	}
	updateKeys := make([]string, 0, len(deployments))
	updateDeployments := make([]rediscache.Value[*cacheentity.DeploymentInfo], 0, len(deployments))
	for _, deployment := range deployments {
		deployment.Cancel = true
		updateKeys = append(updateKeys, repo.formatKey(deployment.ID))
		updateDeployments = append(updateDeployments, rediscache.NewJSONValue(deployment))
	}
	err = rediscache.MSet(ctx, repo.client, updateKeys, updateDeployments, redis.KeepTTL)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (repo *deploymentInfoRepo) Del(ctx context.Context, deploymentID string) error {
	err := rediscache.Del(ctx, repo.client, repo.formatKey(deploymentID))
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (repo *deploymentInfoRepo) formatKey(deploymentID string) string {
	return fmt.Sprintf("deployment:%s:info", deploymentID)
}
