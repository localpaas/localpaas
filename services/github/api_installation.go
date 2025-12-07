package github

import (
	"context"

	"github.com/google/go-github/v75/github"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

func (c *Client) ListInstallations(ctx context.Context, paging *basedto.Paging) (
	[]*github.Installation, *basedto.PagingMeta, error) {
	if !c.isAppClient() {
		return nil, nil, apperrors.Wrap(ErrGithubAppClientRequired)
	}
	opts := &github.ListOptions{
		PerPage: defaultListPageSize,
		Page:    0,
	}
	if paging != nil {
		opts.Page = paging.Offset / gofn.Coalesce(paging.Limit, 1)
		opts.PerPage = paging.Limit
	}
	output, _, err := c.client.Apps.ListInstallations(ctx, opts)
	if err != nil {
		return nil, nil, apperrors.Wrap(err)
	}
	return output, &basedto.PagingMeta{
		Offset: opts.Page * opts.PerPage,
		Limit:  opts.PerPage,
		Total:  -1,
	}, nil
}

func (c *Client) ListAllInstallations(ctx context.Context, options ...ListOption) ([]*github.Installation, error) {
	if !c.isAppClient() {
		return nil, apperrors.Wrap(ErrGithubAppClientRequired)
	}
	output, err := listAll(ctx, c.client,
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
