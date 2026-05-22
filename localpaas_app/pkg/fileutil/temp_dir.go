package fileutil

import (
	"os"
	"path"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
)

// CreateTempDir creates a temp dir.
// Should use "*" for `pattern` value. If empty, only base dir is created. See os.MkdirTemp.
func CreateTempDir(baseDir, pattern string, perm os.FileMode) (dir string, err error) {
	if perm == 0 {
		perm = defaultDirMode
	}
	dir = path.Join(baseDir, timeutil.NowUTC().Format(time.DateOnly))

	err = os.MkdirAll(dir, perm)
	if err != nil {
		return "", apperrors.Wrap(err)
	}

	if pattern != "" {
		dir, err = os.MkdirTemp(dir, pattern)
		if err != nil {
			return "", apperrors.Wrap(err)
		}
	}

	return dir, nil
}
