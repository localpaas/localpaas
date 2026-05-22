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
)

func (s *service) sysCleanupFiles(
	_ context.Context,
	data *sysCleanupData,
) (err error) {
	defer func() {
		if err != nil {
			data.TaskOutput.FileCleanup.Error = err.Error()
		}
	}()

	// Remove outdated temp files
	err1 := s.sysCleanupTempFiles()

	// TODO: add implementation

	return errors.Join(err1)
}

func (s *service) sysCleanupTempFiles() (err error) {
	baseDirs := []string{base.BaseTempDirDefault, filepath.Join(config.Current.AppPath, "tmp")}
	threshold := time.Now().AddDate(0, 0, -3) //nolint:mnd

	for _, baseDir := range baseDirs {
		entries, err := os.ReadDir(baseDir)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return apperrors.Wrap(err)
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			dirTime, err := time.Parse(time.DateOnly, entry.Name())
			if err != nil {
				continue
			}

			if dirTime.Before(threshold) {
				_ = os.RemoveAll(filepath.Join(baseDir, entry.Name()))
			}
		}
	}

	return nil
}
