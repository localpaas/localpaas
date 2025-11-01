package copier

import (
	"github.com/tiendc/go-deepcopy"
)

func Copy(dst, src any) error {
	return deepcopy.Copy(dst, src) //nolint:wrapcheck
}
