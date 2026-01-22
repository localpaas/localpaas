package nginx

import (
	"io"

	crossplane "github.com/localpaas/nginx-go-crossplane"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

type BuildOptions struct {
	crossplane.BuildOptions
}
type BuildOption func(*BuildOptions)

func Build(config *Config, w io.Writer, options ...BuildOption) error {
	if len(config.inner.Parsed) == 0 {
		return nil
	}

	opts := &BuildOptions{}
	for _, opt := range options {
		opt(opts)
	}

	err := crossplane.Build(w, *config.inner, &opts.BuildOptions)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
