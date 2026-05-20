package githelper

import (
	"context"
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

const (
	checkoutMaxDepthDefault = 100
)

func Checkout(
	ctx context.Context,
	checkoutPath string,
	gitOpts *git.CloneOptions,
	commitHash string,
	maxDepth uint,
) (repo *git.Repository, commit *object.Commit, err error) {
	// Try cloning with depth = 2 when commitHash is given
	gitOpts.Depth = gofn.If(commitHash == "", 1, 2) //nolint:mnd

	repo, err = git.PlainCloneContext(ctx, checkoutPath, false, gitOpts)
	if err != nil {
		return nil, nil, apperrors.Wrap(err)
	}

	if commitHash == "" {
		head, err := repo.Head()
		if err != nil {
			return nil, nil, apperrors.Wrap(err)
		}
		commit, err = repo.CommitObject(head.Hash())
		if err != nil {
			return nil, nil, apperrors.Wrap(err)
		}
		return repo, commit, nil
	}

	targetHash := plumbing.NewHash(commitHash)
	depth := uint(0)
	depthIncrement := max(10, maxDepth/10) //nolint:mnd
	maxDepth = gofn.Coalesce(maxDepth, checkoutMaxDepthDefault)
	for depth <= maxDepth {
		commit, err = repo.CommitObject(targetHash)
		if err == nil && commit != nil {
			break
		}
		if !errors.Is(err, plumbing.ErrObjectNotFound) {
			return nil, nil, apperrors.Wrap(err)
		}
		depth += depthIncrement
		err = repo.FetchContext(ctx, &git.FetchOptions{
			RemoteName: git.DefaultRemoteName,
			Depth:      int(depth), //nolint
			Auth:       gitOpts.Auth,
			Force:      true,
		})
		if err != nil {
			return nil, nil, apperrors.Wrap(err)
		}
	}
	if commit == nil {
		return nil, nil, apperrors.Wrap(plumbing.ErrObjectNotFound)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, nil, apperrors.Wrap(err)
	}
	err = worktree.Checkout(&git.CheckoutOptions{
		Hash: targetHash,
	})
	if err != nil {
		return nil, nil, apperrors.Wrap(err)
	}

	return repo, commit, nil
}
