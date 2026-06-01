package syscleanupserviceimpl

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/tasklog"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
)

const (
	repoCacheOutdatedPeriod = 10 * 24 * time.Hour
)

func (s *service) sysCleanupCache(
	ctx context.Context,
	db database.IDB,
	data *sysCleanupData,
) (err error) {
	defer func() {
		if err != nil {
			data.TaskOutput.FileCleanup.Error = err.Error()
		}
	}()

	// Remove old repo cache files in local
	err1 := s.sysCleanupRepoCacheFiles(ctx, db, data)

	// TODO: add more cleanup

	return errors.Join(err1)
}

func (s *service) sysCleanupRepoCacheFiles(
	ctx context.Context,
	db database.IDB,
	data *sysCleanupData,
) (err error) {
	timeNow := timeutil.NowUTC()
	retention := repoCacheOutdatedPeriod

	deletingFiles, _, err := s.fileRepo.List(ctx, db, nil,
		bunex.SelectWhere("file.type = ?", base.FileTypeRepoCache),
		bunex.SelectWhere("file.storage_type = ?", base.FileStorageLocal),
		bunex.SelectWhere("file.updated_at < ?", timeNow.Add(-retention)),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}

	for _, file := range deletingFiles {
		file.DeletedAt = timeNow
		data.TaskOutput.FileCleanup.RepoCacheFilesDeleted++
	}
	err = s.fileRepo.UpsertMulti(ctx, db, deletingFiles, entity.FileUpsertingConflictCols,
		[]string{"deleted_at"})
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Delete real files in local
	rootDir := config.Current.AppPath
	for _, file := range deletingFiles {
		filePath := filepath.Join(file.Path, file.Name)
		filePathAbs := filepath.Join(rootDir, filePath)
		err := os.Remove(filePathAbs)
		if err != nil {
			_ = data.LogStore.Add(ctx, tasklog.NewOutFrame("Failed to remove outdated cache file: "+
				filePath+" with error: "+err.Error(), tasklog.TsNow))
		} else {
			_ = data.LogStore.Add(ctx, tasklog.NewOutFrame("Outdated cache file removed: "+filePath,
				tasklog.TsNow))
		}
	}

	return nil
}
