package syscleanupserviceimpl

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/fileutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/tasklog"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/services/aws/s3"
)

func (s *service) sysCleanupBackups(
	ctx context.Context,
	db database.IDB,
	data *sysCleanupData,
) (err error) {
	if !data.SysCleanupSettings.BackupCleanup.Enabled {
		return nil
	}

	defer func() {
		if err != nil {
			data.TaskOutput.BackupCleanup.Error = err.Error()
		}
	}()

	// Remove old backup files in local
	err1 := s.sysCleanupLocalBackupFiles(ctx, db, data)

	// Remove old backup files in cloud
	err2 := s.sysCleanupCloudBackupFiles(ctx, db, data)

	return errors.Join(err1, err2)
}

func (s *service) sysCleanupLocalBackupFiles(
	ctx context.Context,
	db database.IDB,
	data *sysCleanupData,
) (err error) {
	if data.SysCleanupSettings.BackupCleanup.LocalBackupRetention < 0 { // No cleanup
		return nil
	}

	timeNow := timeutil.NowUTC()
	retention := data.SysCleanupSettings.BackupCleanup.LocalBackupRetention.ToDuration()

	deletingFileSettings, _, err := s.settingRepo.List(ctx, db, base.NewSettingScopeGlobal(), nil,
		bunex.SelectWhere("setting.type = ?", base.SettingTypeFile),
		bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
		bunex.SelectWhere("setting.kind = ?", base.FileKindSystemBackup),
		bunex.SelectWhere("setting.data->>'storageType' = ?", base.FileStorageLocal),
		bunex.SelectWhere("setting.created_at < ?", timeNow.Add(-retention)),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}

	for _, setting := range deletingFileSettings {
		setting.DeletedAt = timeNow
	}
	err = s.settingRepo.UpsertMulti(ctx, db, deletingFileSettings,
		entity.SettingUpsertingConflictCols, []string{"deleted_at"})
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Delete real files in local
	backupSaveDir := config.Current.DataPathSystemBackupFiles()
	exists, err := fileutil.FileExists(backupSaveDir, false)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if !exists {
		return nil
	}

	for _, setting := range deletingFileSettings {
		file := setting.MustAsFile()
		filePath := filepath.Join(file.Path, file.Name)
		err := os.Remove(filePath)
		if err != nil {
			_ = data.LogStore.Add(ctx, tasklog.NewOutFrame("Failed to remove outdated backup file: "+
				filePath+" with error: "+err.Error(), tasklog.TsNow))
		} else {
			_ = data.LogStore.Add(ctx, tasklog.NewOutFrame("Outdated backup file removed: "+filePath,
				tasklog.TsNow))
		}
	}

	return nil
}

func (s *service) sysCleanupCloudBackupFiles(
	ctx context.Context,
	db database.IDB,
	data *sysCleanupData,
) (err error) {
	if data.SysCleanupSettings.BackupCleanup.CloudBackupRetention < 0 { // No cleanup
		return nil
	}

	timeNow := timeutil.NowUTC()
	retention := data.SysCleanupSettings.BackupCleanup.CloudBackupRetention.ToDuration()

	deletingFileSettings, _, err := s.settingRepo.List(ctx, db, base.NewSettingScopeGlobal(), nil,
		bunex.SelectWhere("setting.type = ?", base.SettingTypeFile),
		bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
		bunex.SelectWhere("setting.kind = ?", base.FileKindSystemBackup),
		bunex.SelectWhere("setting.data->>'storageType' = ?", base.FileStorageCloud),
		bunex.SelectWhere("setting.created_at < ?", timeNow.Add(-retention)),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}

	for _, setting := range deletingFileSettings {
		setting.DeletedAt = timeNow
	}
	err = s.settingRepo.UpsertMulti(ctx, db, deletingFileSettings,
		entity.SettingUpsertingConflictCols, []string{"deleted_at"})
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Load all reference objects from the files
	refObjects, err := s.settingService.LoadReferenceObjects(ctx, db, nil, true,
		false, deletingFileSettings...)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.AddRefObjects(refObjects)

	// Delete real files in cloud
	mapDelFuncByStorage := map[string]func(*entity.File) error{}

	getDelFunc := func(file *entity.File) (func(*entity.File) error, error) {
		delFunc, exists := mapDelFuncByStorage[file.Storage.ID]
		if exists {
			return delFunc, nil
		}

		storageSttg := data.RefObjects.RefSettings[file.Storage.ID]
		if storageSttg == nil {
			return nil, apperrors.NewNotFound(fmt.Sprintf("Storage setting '%v'", file.Storage.ID))
		}

		switch base.CloudStorageKind(storageSttg.Kind) { //nolint:gocritic
		case base.CloudStorageKindS3:
			s3Client, err := s3.NewClientFromSetting(ctx, storageSttg)
			if err != nil {
				return nil, apperrors.Wrap(err)
			}
			delFunc = func(file *entity.File) error {
				return s3Client.DeleteObject(ctx, file.Bucket, filepath.Join(file.Path, file.Name))
			}
			mapDelFuncByStorage[storageSttg.ID] = delFunc
		}
		return delFunc, nil
	}

	for _, setting := range deletingFileSettings {
		file := setting.MustAsFile()
		filePath := filepath.Join(file.Path, file.Name)

		delFunc, err := getDelFunc(file)
		if err != nil {
			_ = data.LogStore.Add(ctx, tasklog.NewOutFrame("Failed to remove backup file in cloud: "+
				filePath+" with creating client error: "+err.Error(), tasklog.TsNow))
		}

		err = delFunc(file)
		if err != nil {
			_ = data.LogStore.Add(ctx, tasklog.NewOutFrame("Failed to remove backup file in cloud: "+
				filePath+" with error: "+err.Error(), tasklog.TsNow))
		} else {
			_ = data.LogStore.Add(ctx, tasklog.NewOutFrame("Outdated backup file removed from cloud: "+filePath,
				tasklog.TsNow))
		}
	}
	return nil
}
