package taskcronjobexec

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/applog"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
)

func (e *Executor) sysDBBackupRemoveOldFiles(
	ctx context.Context,
	sysBackup *entity.SystemBackup,
	data *sysBackupTaskData,
) (err error) {
	if sysBackup.LocalBackupRetention <= 0 {
		return nil
	}

	bakFiles, err := os.ReadDir(data.BackupFileDir)
	if err != nil {
		return apperrors.Wrap(err)
	}

	oldestTime := timeutil.NowUTC().Add(-sysBackup.LocalBackupRetention.ToDuration())
	for _, entry := range bakFiles {
		if entry.IsDir() {
			continue
		}
		filename := entry.Name()
		fileTime := sysDBBackupParseFileTime(filename)
		if !fileTime.IsZero() && fileTime.Before(oldestTime) {
			err := os.Remove(filepath.Join(data.BackupFileDir, filename))
			if err != nil {
				_ = data.LogStore.Add(ctx, applog.NewOutFrame("Failed to remove outdated backup file: "+
					filename+" with error: "+err.Error(), applog.TsNow))
			} else {
				_ = data.LogStore.Add(ctx, applog.NewOutFrame("Outdated backup file removed: "+filename,
					applog.TsNow))
			}
		}
	}

	return nil
}

func sysDBBackupParseFileTime(filename string) time.Time {
	filename = strings.TrimPrefix(filename, backupFilePrefix)
	timeStr, _, _ := strings.Cut(filename, ".")
	dt, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}
	}
	return dt
}
