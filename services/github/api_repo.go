package github

import (
	"context"

	"github.com/google/go-github/v75/github"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

func (app *App) ListRepos(ctx context.Context, options ...ListOption) ([]*github.Repository, error) {
	opts := &github.ListOptions{}
	for _, option := range options {
		option(opts)
	}
	output, _, err := app.client.Apps.ListRepos(ctx, opts)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return output.Repositories, nil
}

func (app *App) ListAllRepos(ctx context.Context, options ...ListOption) ([]*github.Repository, error) {
	output, err := listAll(ctx, app.client,
		func(ctx context.Context, client *github.Client, opts *github.ListOptions) (
			[]*github.Repository, *github.Response, error) {
			for _, option := range options {
				option(opts)
			}
			output, resp, err := client.Apps.ListRepos(ctx, opts)
			if err != nil {
				return nil, nil, apperrors.Wrap(err)
			}
			return output.Repositories, resp, nil
		})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return output, nil
}
