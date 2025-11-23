package nginx

import (
	"io"
	"strings"

	crossplane "github.com/nginxinc/nginx-go-crossplane"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

type ParseOption func(*crossplane.ParseOptions)

func ParseFile(filepath string, options ...ParseOption) (*Config, error) {
	opts := &crossplane.ParseOptions{}
	for _, opt := range options {
		opt(opts)
	}

	payload, err := crossplane.Parse(filepath, opts)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp := &Config{}
	if len(payload.Config) > 0 {
		resp.inner = &payload.Config[0]
	}
	return resp, nil
}

func ParseString(data string, options ...ParseOption) (*Config, error) {
	opts := &crossplane.ParseOptions{
		SingleFile: true,
		Open: func(path string) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader(data)), nil
		},
	}
	for _, opt := range options {
		opt(opts)
	}

	payload, err := crossplane.Parse("nginx.conf", opts)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp := &Config{}
	if len(payload.Config) > 0 {
		resp.inner = &payload.Config[0]
	}
	return resp, nil
}
