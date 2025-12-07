package github

import (
	"context"

	"github.com/google/go-github/v75/github"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

func (c *Client) ListRepos(ctx context.Context, options ...ListOption) ([]*github.Repository, error) {
	if c.isAppClient() {
		return c.listAppRepos(ctx, options...)
	}
	return c.listUserRepos(ctx, options...)
}

func (c *Client) listAppRepos(ctx context.Context, options ...ListOption) ([]*github.Repository, error) {
	opts := &github.ListOptions{}
	for _, option := range options {
		option(opts)
	}
	output, _, err := c.client.Apps.ListRepos(ctx, opts)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return output.Repositories, nil
}

func (c *Client) listUserRepos(ctx context.Context, options ...ListOption) ([]*github.Repository, error) {
	opts := &github.ListOptions{}
	for _, option := range options {
		option(opts)
	}
	output, _, err := c.client.Repositories.ListByAuthenticatedUser(ctx,
		&github.RepositoryListByAuthenticatedUserOptions{
			ListOptions: *opts,
		})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return output, nil
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
