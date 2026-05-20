package githelper

import (
	"errors"

	"github.com/go-git/go-git/v5/plumbing"
)

func IsErrObjectNotFound(err error) bool {
	return errors.Is(err, plumbing.ErrObjectNotFound)
}
