package tracerr

import goerrors "github.com/go-errors/errors"

// Wrap wraps an error with adding stack trace
func Wrap(err error) error {
	if err == nil {
		return nil
	}
	return goerrors.Wrap(err, 1)
}
