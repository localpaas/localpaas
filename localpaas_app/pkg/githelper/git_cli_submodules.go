package githelper

import (
	"context"
	"os/exec"

	"github.com/go-git/go-git/v5"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
)

func gitCliFetchSubmodules(
	ctx context.Context,
	checkoutOpts *CheckoutOptions,
) (err error) {
	if checkoutOpts.RecurseSubmodules == git.NoRecurseSubmodules {
		return nil
	}
	args := []string{"submodule", "update", "--init", "--recursive"}
	if checkoutOpts.ShallowSubmodules {
		args = append(args, "--depth", "1")
	}

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = checkoutOpts.CheckoutDir
	cmd.Env = checkoutOpts.sharedEnv

	out, err := cmd.CombinedOutput()
	addLog(ctx, reflectutil.UnsafeBytesToStr(out), err != nil, checkoutOpts)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
