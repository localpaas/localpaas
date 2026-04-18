package fileutil

import (
	"path/filepath"
	"strings"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

func IsSubpath(basepath, targetpath string) (bool, error) {
	// Rel returns a relative path that is lexically equivalent to targetpath when joined to basepath
	rel, err := filepath.Rel(basepath, targetpath)
	if err != nil {
		return false, apperrors.Wrap(err)
	}
	// Two paths are the same
	if rel == "." {
		return false, nil
	}
	// A result starting with `..` means the target is outside the base directory
	if strings.HasPrefix(rel, "..") {
		return false, nil
	}
	return true, nil
}

func IsEqualOrSubpath(basepath, targetpath string) (bool, error) {
	// Rel returns a relative path that is lexically equivalent to targetpath when joined to basepath
	rel, err := filepath.Rel(basepath, targetpath)
	if err != nil {
		return false, apperrors.Wrap(err)
	}
	if rel == "." {
		return true, nil
	}
	// A result starting with `..` means the target is outside the base directory
	if strings.HasPrefix(rel, "..") {
		return false, nil
	}
	return true, nil
}
