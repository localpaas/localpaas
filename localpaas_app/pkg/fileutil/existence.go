package fileutil

import (
	"errors"
	"os"

	"github.com/localpaas/localpaas/localpaas_app/pkg/tracerr"
)

func FileExists(filename string, isFile bool) (bool, error) {
	fileInfo, err := os.Stat(filename)
	if err == nil && fileInfo.IsDir() == !isFile {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, tracerr.Wrap(err)
}
