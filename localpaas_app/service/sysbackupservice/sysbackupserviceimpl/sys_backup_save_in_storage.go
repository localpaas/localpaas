package sysbackupserviceimpl

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/tasklog"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/services/aws/s3"
)

func (s *service) sysBackupSaveResultInStorage(
	ctx context.Context,
	db database.IDB,
	data *sysBackupData,
) (err error) {
	if data.SysBackupSettings.CloudStorage.ID == "" {
		return nil
	}
	storageSttg := data.RefObjects.RefSettings[data.SysBackupSettings.CloudStorage.ID]
	if storageSttg == nil {
		return nil
	}

	var storageName string
	var storageBucket string
	var uploadFunc func(targetFilePath string, data io.Reader) error

	switch base.CloudStorageKind(storageSttg.Kind) {
	case base.CloudStorageKindS3:
		s3Client, err := s3.NewClientFromSetting(ctx, storageSttg)
		if err != nil {
			return apperrors.Wrap(err)
		}
		storageName = "AWS S3"
		storageBucket = s3Client.Config.Bucket
		uploadFunc = func(targetFilePath string, input io.Reader) error {
			return s3Client.UploadEx(ctx, storageBucket, targetFilePath,
				0, 0, input)
		}
	default:
		return apperrors.NewUnsupported(fmt.Sprintf("Storage type '%v'", storageSttg.Kind))
	}

	targetFilePath := filepath.Join(data.SysBackupSettings.CloudStorage.DestinationDir, data.OutFileName)
	backupFile, err := os.Open(data.OutFilePath)
	if err != nil {
		return apperrors.Wrap(err)
	}
	defer backupFile.Close()

	start := timeutil.NowUTC()
	_ = data.LogStore.Add(ctx, tasklog.NewOutFrame(fmt.Sprintf(
		"Start uploading file '%v' to '%v' bucket '%v'...",
		data.OutFileName, storageName, storageBucket), tasklog.TsNow))

	err = uploadFunc(targetFilePath, backupFile)
	if err != nil {
		_ = data.LogStore.Add(ctx, tasklog.NewWarnFrame(
			"Failed to upload backup file to "+storageName+" with error: "+err.Error(), tasklog.TsNow))
		return apperrors.Wrap(err)
	}
	_ = data.LogStore.Add(ctx, tasklog.NewOutFrame("Backup file uploaded to "+storageName+
		" in "+time.Since(start).String(), tasklog.TsNow))

	localFile := data.LocalOutFile.MustAsFile()
	remoteFileSetting := &entity.Setting{
		ID:        gofn.Must(ulid.NewStringULID()),
		Scope:     base.SettingScopeGlobal,
		Type:      base.SettingTypeFile,
		Kind:      string(base.FileKindSystemBackup),
		Status:    base.SettingStatusActive,
		Name:      data.OutFileName,
		Version:   entity.CurrentFileVersion,
		CreatedAt: data.TimeNow,
		UpdatedAt: data.TimeNow,
	}

	remoteFile := &entity.File{
		FileKind:    base.FileKindSystemBackup,
		StorageType: base.FileStorageCloud,
		Storage:     entity.ObjectID{ID: data.SysBackupSettings.CloudStorage.ID},
		Bucket:      storageBucket,
		Mimetype:    localFile.Mimetype,
		Name:        data.OutFileName,
		Size:        localFile.Size,
		Path:        data.SysBackupSettings.CloudStorage.DestinationDir,
	}

	remoteFileSetting.MustSetData(remoteFile)
	data.RemoteOutFile = remoteFileSetting

	err = s.settingRepo.Insert(ctx, db, remoteFileSetting)
	if err != nil {
		_ = data.LogStore.Add(ctx, tasklog.NewOutFrame("Failed to save file record into DB with error: "+
			err.Error(), tasklog.TsNow))
		return apperrors.Wrap(err)
	}

	return nil
}
