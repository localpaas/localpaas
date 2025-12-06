package github

import (
	"context"

	"github.com/google/go-github/v75/github"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

func (app *App) ListInstallations(ctx context.Context, paging *basedto.Paging) (
	[]*github.Installation, *basedto.PagingMeta, error) {
	opts := &github.ListOptions{
		PerPage: defaultPerPage,
		Page:    0,
	}
	if paging != nil {
		opts.Page = paging.Offset / paging.Limit
		opts.PerPage = paging.Limit
	}
	output, _, err := app.client.Apps.ListInstallations(ctx, opts)
	if err != nil {
		return nil, nil, apperrors.Wrap(err)
	}
	return output, &basedto.PagingMeta{
		Offset: opts.Page * opts.PerPage,
		Limit:  opts.PerPage,
		Total:  -1,
	}, nil
}

func (app *App) ListAllInstallations(ctx context.Context, options ...ListOption) ([]*github.Installation, error) {
	output, err := listAll(ctx, app.client,
		func(ctx context.Context, client *github.Client, opts *github.ListOptions) (
			[]*github.Installation, *github.Response, error) {
			for _, option := range options {
				option(opts)
			}
			return client.Apps.ListInstallations(ctx, opts)
		})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return output, nil
}
