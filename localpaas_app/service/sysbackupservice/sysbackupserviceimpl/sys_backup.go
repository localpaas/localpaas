package sysbackupserviceimpl

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"filippo.io/age"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/funcutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/jsonl"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/service/sysbackupservice"
)

const (
	sysBackupVer = "v1.0.0"

	// 0755 grants read/write/execute for owner, read/execute for group/others
	// 0644 grants read/write for owner, read-only for group/others
	sysBackupDirFileMode = 0o755
)

type sysBackupData struct {
	*sysbackupservice.SysBackupReq
	TaskOutput *entity.TaskSystemBackupOutput
	TimeNow    time.Time

	BackupRootDir string
	BackupSaveDir string

	OutFileName   string
	OutFilePath   string
	LocalOutFile  *entity.Setting
	RemoteOutFile *entity.Setting
}

func (s *service) Backup(
	ctx context.Context,
	db database.Tx,
	req *sysbackupservice.SysBackupReq,
) (resp *sysbackupservice.SysBackupResp, err error) {
	defer funcutil.EnsureNoPanic(&err)

	resp = &sysbackupservice.SysBackupResp{}
	data := &sysBackupData{
		SysBackupReq: req,
		TaskOutput: &entity.TaskSystemBackupOutput{
			DBBackup: &entity.DBBackupOutput{},
		},
		TimeNow: timeutil.NowUTC(),
	}

	// Backup DB
	err = s.sysBackup(ctx, db, data)

	// Assign back the result output
	data.Task.MustSetOutput(data.TaskOutput)

	return resp, nil
}

func (s *service) sysBackup(
	ctx context.Context,
	db database.IDB,
	data *sysBackupData,
) (err error) {
	defer func() {
		if err != nil {
			data.TaskOutput.DBBackup.Error = err.Error()
		}
	}()

	tmpFile, jsonlW, closer, err := s.sysBackupCreateWriter(data)
	if err != nil {
		return apperrors.Wrap(err)
	}
	defer func() {
		if closer != nil {
			_ = closer(true)
		}
	}()

	// Write header to the backup file
	err = jsonlW.WriteMetadata(jsonl.Metadata{
		Type:      "system-backup",
		Version:   sysBackupVer,
		Timestamp: data.TimeNow.Truncate(time.Second),
	})
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Start the data backup
	err = s.sysBackupDB(ctx, db, jsonlW, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = s.sysBackupFiles(ctx, db, jsonlW, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	_ = closer(false) // Flush data in writers, but not remove the temp file
	closer = nil

	// Save the result in a local file
	err = s.sysBackupSaveResultInLocal(ctx, db, tmpFile, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Upload backup file to cloud storage if configured
	err = s.sysBackupSaveResultInStorage(ctx, db, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (s *service) sysBackupCreateWriter(
	data *sysBackupData,
) (tmpFile *os.File, jsonlW *jsonl.Writer, closer func(bool) error, err error) {
	// Make sure the backup directory exist
	data.BackupRootDir = config.Current.DataPathSystemBackup()
	data.BackupSaveDir = config.Current.DataPathSystemBackupFiles()

	err = os.MkdirAll(data.BackupRootDir, sysBackupDirFileMode)
	if err != nil {
		return nil, nil, nil, apperrors.Wrap(err)
	}

	tmpDir := filepath.Join(data.BackupRootDir, "tmp")
	err = os.MkdirAll(tmpDir, sysBackupDirFileMode)
	if err != nil {
		return nil, nil, nil, apperrors.Wrap(err)
	}

	tmpFile, err = os.CreateTemp(tmpDir, "*.bak")
	if err != nil {
		return nil, nil, nil, apperrors.Wrap(err)
	}

	defer func() {
		if err != nil {
			_ = os.Remove(tmpFile.Name())
		}
	}()

	var w io.Writer
	var encW, gzW io.WriteCloser
	w = tmpFile

	switch data.SysBackupSettings.Encryption.Format {
	case base.FileEncryptionFormatAge:
		encSecret := data.SysBackupSettings.Encryption.Secret.MustGetPlain()
		if encSecret == "" {
			return nil, nil, nil, apperrors.NewMissing("Encryption secret")
		}
		recipient, err := age.NewScryptRecipient(encSecret)
		if err != nil {
			return nil, nil, nil, apperrors.Wrap(err)
		}
		encW, err = age.Encrypt(w, recipient)
		if err != nil {
			return nil, nil, nil, apperrors.Wrap(err)
		}
		w = encW
	case base.FileEncryptionNone: // Do nothing
	default:
		return nil, nil, nil, apperrors.NewUnsupported(
			fmt.Sprintf("Encryption format '%v'", data.SysBackupSettings.Encryption.Format))
	}

	switch data.SysBackupSettings.Compression.Format {
	case base.FileCompressionFormatGzip:
		gzW = gzip.NewWriter(w)
		w = gzW
	case base.FileCompressionNone: // Do nothing
	default:
		return nil, nil, nil, apperrors.NewUnsupported(
			fmt.Sprintf("Compression format '%v'", data.SysBackupSettings.Compression.Format))
	}

	jsonlW = jsonl.NewWriter(w)

	closer = func(removeTmpFile bool) (err error) {
		if e := jsonlW.Close(); e != nil {
			err = errors.Join(err, e)
		}
		if gzW != nil {
			if e := gzW.Close(); e != nil {
				err = errors.Join(err, e)
			}
		}
		if encW != nil {
			if e := encW.Close(); e != nil {
				err = errors.Join(err, e)
			}
		}
		if e := tmpFile.Close(); e != nil {
			err = errors.Join(err, e)
		}
		if removeTmpFile {
			_ = os.Remove(tmpFile.Name()) // Ignore this error as the temp file may not exist
		}
		return err
	}

	return tmpFile, jsonlW, closer, nil
}
