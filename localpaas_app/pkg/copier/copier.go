package copier

import (
	"github.com/tiendc/go-deepcopy"
)

func Copy(dst, src any) error {
	return deepcopy.Copy(dst, src) //nolint:wrapcheck
}

func CopySlice[T any](src []T) ([]T, error) {
	var dst []T
	if err := deepcopy.Copy(&dst, &src); err != nil {
		return nil, err //nolint:wrapcheck
	}
	return dst, nil
}
