package nginx

import (
	"io"

	crossplane "github.com/nginxinc/nginx-go-crossplane"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

type BuildOption func(*crossplane.BuildOptions)

func Build(config *Config, w io.Writer, options ...BuildOption) error {
	if len(config.inner.Parsed) == 0 {
		return nil
	}

	opts := &crossplane.BuildOptions{}
	for _, opt := range options {
		opt(opts)
	}

	err := crossplane.Build(w, *config.inner, opts)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
