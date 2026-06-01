package githelper

import (
	"context"
	"os"
	"os/exec"
	"strconv"

	"github.com/go-git/go-git/v5"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
)

func gitCliClone(
	ctx context.Context,
	checkoutOpts *CheckoutOptions,
) (repo *git.Repository, err error) {
	err = os.MkdirAll(checkoutOpts.CheckoutDir, base.DirModeDefault)
	if err != nil {
		return nil, apperrors.New(err)
	}

	args := []string{"clone"}
	if checkoutOpts.SingleBranch {
		args = append(args, "--single-branch")
	}
	if checkoutOpts.Depth > 0 {
		args = append(args, "--depth", strconv.Itoa(int(checkoutOpts.Depth))) //nolint:gosec
	}
	if checkoutOpts.branch != "" {
		args = append(args, "--branch", checkoutOpts.branch)
	}
	args = append(args, "--", checkoutOpts.URL, checkoutOpts.CheckoutDir)

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Env = checkoutOpts.sharedEnv
	out, err := cmd.CombinedOutput()
	addLog(ctx, reflectutil.UnsafeBytesToStr(out), err != nil, checkoutOpts)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	// Open repo with go-git
	repo, err = git.PlainOpen(checkoutOpts.CheckoutDir)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return repo, nil
}
