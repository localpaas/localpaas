package imagebuildsettingsuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/service/syscleanupservice"
	"github.com/localpaas/localpaas/localpaas_app/tasks/queue"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/imagebuildsettingsuc/imagebuildsettingsdto"
)

func (uc *UC) ClearRepoCache(
	ctx context.Context,
	auth *basedto.Auth,
	req *imagebuildsettingsdto.ClearRepoCacheReq,
) (*imagebuildsettingsdto.ClearRepoCacheResp, error) {
	cleanupReq := &syscleanupservice.SysCleanupReq{
		TaskExecData: &queue.TaskExecData{
			Task: &entity.Task{},
		},
		SysCleanupSettings: &entity.SystemCleanup{
			CacheCleanup: entity.SystemCacheCleanup{
				Enabled: true,
			},
		},
		CleanupCacheRepo: syscleanupservice.CleanupFlagForce,
	}

	filesDeleted := 0
	spaceReclaimed := uint64(0)
	err := transaction.Execute(ctx, uc.DB, func(db database.Tx) error {
		resp, err := uc.sysCleanupService.Cleanup(ctx, db, cleanupReq)
		if err != nil {
			return apperrors.New(err)
		}
		filesDeleted = resp.TaskOutput.CacheCleanup.RepoCacheFilesDeleted
		spaceReclaimed = resp.TaskOutput.CacheCleanup.RepoCacheSpaceReclaimed
		return nil
	})
	if err != nil {
		return nil, apperrors.New(err)
	}

	return &imagebuildsettingsdto.ClearRepoCacheResp{
		Data: &imagebuildsettingsdto.ClearRepoCacheDataResp{
			FilesDeleted:   filesDeleted,
			SpaceReclaimed: spaceReclaimed,
		},
	}, nil
}
