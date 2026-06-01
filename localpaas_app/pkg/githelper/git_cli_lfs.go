package githelper

import (
	"context"
	"os/exec"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
)

//nolint:unused
func gitCliPullLfs(
	ctx context.Context,
	checkoutOpts *CheckoutOptions,
) (err error) {
	if !checkoutOpts.LFSEnabled {
		return nil
	}

	cmd := exec.CommandContext(ctx, "git", "lfs", "pull")
	cmd.Dir = checkoutOpts.CheckoutDir
	cmd.Env = checkoutOpts.sharedEnv

	out, err := cmd.CombinedOutput()
	addLog(ctx, reflectutil.UnsafeBytesToStr(out), err != nil, checkoutOpts)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
