package nginx

import "errors"

var (
	ErrServerBlockRequired = errors.New("server block is required")
)
