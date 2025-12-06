package github

import (
	"context"

	"github.com/google/go-github/v75/github"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

const (
	defaultPerPage = 100
)

type ListOption func(*github.ListOptions)

func ListOptionPage(p int) ListOption {
	return func(opts *github.ListOptions) {
		opts.Page = p
	}
}

func ListOptionPerPage(p int) ListOption {
	return func(opts *github.ListOptions) {
		opts.PerPage = p
	}
}

type listFunc[T any] func(context.Context, *github.Client, *github.ListOptions) ([]T, *github.Response, error)

func listAll[T any](ctx context.Context, client *github.Client, fn listFunc[T]) (output []T, err error) {
	opts := &github.ListOptions{
		Page:    0,
		PerPage: defaultPerPage,
	}
	for {
		result, resp, err := fn(ctx, client, opts)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		output = append(output, result...)
		if resp.NextPage == 0 {
			break
		}
		if resp.Rate.Remaining == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return output, nil
}
