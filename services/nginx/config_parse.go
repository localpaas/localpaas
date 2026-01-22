package nginx

import (
	"io"
	"strings"

	crossplane "github.com/localpaas/nginx-go-crossplane"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

type ParseOptions struct {
	crossplane.ParseOptions
}

type ParseOption func(*ParseOptions)

func ParseFile(filepath string, options ...ParseOption) (*Config, error) {
	opts := &ParseOptions{}
	for _, opt := range options {
		opt(opts)
	}

	payload, err := crossplane.Parse(filepath, &opts.ParseOptions)
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
	opts := &ParseOptions{
		ParseOptions: crossplane.ParseOptions{
			SingleFile: true,
			Open: func(path string) (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader(data)), nil
			},
		},
	}
	for _, opt := range options {
		opt(opts)
	}

	payload, err := crossplane.Parse("nginx.conf", &opts.ParseOptions)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp := &Config{}
	if len(payload.Config) > 0 {
		resp.inner = &payload.Config[0]
	}
	return resp, nil
}
