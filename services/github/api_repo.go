package github

import (
	"context"

	"github.com/google/go-github/v75/github"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

func (c *Client) ListRepos(ctx context.Context, paging *basedto.Paging) (
	[]*github.Repository, *basedto.PagingMeta, error) {
	if c.isAppClient() {
		return c.listAppRepos(ctx, paging)
	}
	return c.listUserRepos(ctx, paging)
}

func (c *Client) listAppRepos(ctx context.Context, paging *basedto.Paging) (
	[]*github.Repository, *basedto.PagingMeta, error) {
	opts := &github.ListOptions{
		PerPage: defaultListPageSize,
		Page:    0,
	}
	if paging != nil {
		opts.Page = paging.Offset / gofn.Coalesce(paging.Limit, 1)
		opts.PerPage = paging.Limit
	}
	output, _, err := c.client.Apps.ListRepos(ctx, opts)
	if err != nil {
		return nil, nil, apperrors.Wrap(err)
	}
	return output.Repositories, &basedto.PagingMeta{
		Offset: opts.Page * opts.PerPage,
		Limit:  opts.PerPage,
		Total:  -1,
	}, nil
}

func (c *Client) listUserRepos(ctx context.Context, paging *basedto.Paging) (
	[]*github.Repository, *basedto.PagingMeta, error) {
	opts := &github.ListOptions{
		PerPage: defaultListPageSize,
		Page:    0,
	}
	if paging != nil {
		opts.Page = paging.Offset / gofn.Coalesce(paging.Limit, 1)
		opts.PerPage = paging.Limit
	}
	output, _, err := c.client.Repositories.ListByAuthenticatedUser(ctx,
		&github.RepositoryListByAuthenticatedUserOptions{
			ListOptions: *opts,
		})
	if err != nil {
		return nil, nil, apperrors.Wrap(err)
	}
	return output, &basedto.PagingMeta{
		Offset: opts.Page * opts.PerPage,
		Limit:  opts.PerPage,
		Total:  -1,
	}, nil
}

func (c *Client) ListAllRepos(ctx context.Context, options ...ListOption) ([]*github.Repository, error) {
	if c.isAppClient() {
		return c.listAllAppRepos(ctx, options...)
	}
	return c.listAllUserRepos(ctx, options...)
}

func (c *Client) listAllAppRepos(ctx context.Context, options ...ListOption) ([]*github.Repository, error) {
	output, err := listAll(ctx, c.client,
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

func (c *Client) listAllUserRepos(ctx context.Context, options ...ListOption) ([]*github.Repository, error) {
	listOpts := &github.RepositoryListByAuthenticatedUserOptions{}
	output, err := listAll(ctx, c.client,
		func(ctx context.Context, client *github.Client, opts *github.ListOptions) (
			[]*github.Repository, *github.Response, error) {
			for _, option := range options {
				option(opts)
			}
			listOpts.ListOptions = *opts
			output, resp, err := c.client.Repositories.ListByAuthenticatedUser(ctx, listOpts)
			if err != nil {
				return nil, nil, apperrors.Wrap(err)
			}
			return output, resp, nil
		})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return output, nil
}
