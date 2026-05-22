package githelper

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gitsight/go-vcsurl"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

const (
	checkoutMaxDepthDefault = 100
	checkoutPathFileMode    = 0755
	sshKeyFileMode          = 0600
)

type AuthSSHKey struct {
	*ssh.PublicKeys
	PEMBytes []byte
}

type CheckoutOptions struct {
	URL               string
	Auth              transport.AuthMethod
	RemoteName        string
	ReferenceName     plumbing.ReferenceName
	SingleBranch      bool
	Depth             uint
	MaxDepth          uint
	RecurseSubmodules git.SubmoduleRescursivity
	ShallowSubmodules bool
	CommitHash        string

	TempDir      string
	CheckoutPath string

	branch string
	sshCmd string
}

func CheckoutWithGitCli(
	ctx context.Context,
	checkoutOpts *CheckoutOptions,
) (repo *git.Repository, commit *object.Commit, err error) {
	// 1. Prepare args
	err = gitCliProcessCheckoutOpts(checkoutOpts)
	if err != nil {
		return nil, nil, apperrors.New(err)
	}

	// 2. Clone repository using git cli
	repo, err = gitCliClone(ctx, checkoutOpts)
	if err != nil {
		return nil, nil, apperrors.Wrap(err)
	}

	// 3. Checkout target commit
	commit, err = gitCliCheckoutTargetCommit(ctx, repo, checkoutOpts)
	if err != nil {
		return nil, nil, apperrors.Wrap(err)
	}

	// 4. Fetch submodules if needed
	err = gitCliFetchSubmodules(ctx, checkoutOpts)
	if err != nil {
		return nil, nil, apperrors.Wrap(err)
	}

	return repo, commit, nil
}

func gitCliProcessCheckoutOpts(
	checkoutOpts *CheckoutOptions,
) (err error) {
	if checkoutOpts.Auth != nil { //nolint:nestif
		parseURL, err := vcsurl.Parse(checkoutOpts.URL)
		if err != nil {
			return apperrors.Wrap(err)
		}

		switch auth := checkoutOpts.Auth.(type) {
		case *http.BasicAuth:
			// Use https url
			if !strings.HasPrefix(checkoutOpts.URL, "https://") {
				checkoutOpts.URL = GetHttpsUrl(parseURL)
			}
			// Add user info to the url
			u, err := url.Parse(checkoutOpts.URL)
			if err != nil {
				return apperrors.Wrap(err)
			}
			u.User = url.UserPassword(auth.Username, auth.Password)
			checkoutOpts.URL = u.String()

		case *AuthSSHKey:
			// Use SSH key to clone, the url must be like `git@host.domain:user/repo.git`
			if !strings.HasPrefix(checkoutOpts.URL, "git@") {
				checkoutOpts.URL = GetSshUrl(parseURL)
			}

			sshKeyFile, err := writeSshKeyFile(checkoutOpts.TempDir, auth.PEMBytes)
			if err != nil {
				return apperrors.Wrap(err)
			}
			checkoutOpts.sshCmd = "ssh -o StrictHostKeyChecking=no -i " + sshKeyFile

		default:
			return apperrors.NewUnsupported(fmt.Sprintf("Git auth method '%v'", auth.Name()))
		}
	}

	if checkoutOpts.Depth == 0 {
		checkoutOpts.Depth = 1
	}
	checkoutOpts.branch = checkoutOpts.ReferenceName.Short()

	return nil
}

func writeSshKeyFile(baseDir string, pemBytes []byte) (sshKeyFile string, err error) {
	fh, err := os.CreateTemp(baseDir, "git-ssh-*")
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	defer fh.Close()

	// NOTE: file will be removed along with the whole temp dir by the caller
	sshKeyFile = fh.Name()

	if err := os.Chmod(sshKeyFile, sshKeyFileMode); err != nil {
		return "", apperrors.Wrap(err)
	}

	if _, err := fh.Write(pemBytes); err != nil {
		return "", apperrors.Wrap(err)
	}

	if pemBytes[len(pemBytes)-1] != '\n' {
		if _, err := fh.Write([]byte("\n")); err != nil {
			return "", apperrors.Wrap(err)
		}
	}

	return sshKeyFile, nil
}

func gitCliClone(
	ctx context.Context,
	checkoutOpts *CheckoutOptions,
) (repo *git.Repository, err error) {
	err = os.MkdirAll(checkoutOpts.CheckoutPath, checkoutPathFileMode)
	if err != nil {
		return nil, apperrors.New(err)
	}

	args := []string{"clone"}
	if checkoutOpts.SingleBranch {
		args = append(args, "--single-branch")
	}
	if checkoutOpts.Depth > 0 {
		args = append(args, "--depth", strconv.FormatUint(uint64(checkoutOpts.Depth), 10))
	}
	if checkoutOpts.branch != "" {
		args = append(args, "--branch", checkoutOpts.branch)
	}
	args = append(args, "--", checkoutOpts.URL, checkoutOpts.CheckoutPath)

	env := []string{} // No use current process's env
	if checkoutOpts.sshCmd != "" {
		env = append(env, "GIT_SSH_COMMAND="+checkoutOpts.sshCmd)
	}

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Env = env
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, apperrors.New(err).WithExtraDetail("%s", string(out))
	}

	// Open repo with go-git
	repo, err = git.PlainOpen(checkoutOpts.CheckoutPath)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return repo, nil
}

func gitCliCheckoutTargetCommit(
	ctx context.Context,
	repo *git.Repository,
	checkoutOpts *CheckoutOptions,
) (commit *object.Commit, err error) {
	if checkoutOpts.CommitHash == "" {
		head, err := repo.Head()
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		commit, err = repo.CommitObject(head.Hash())
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		return commit, nil
	}

	targetHash := plumbing.NewHash(checkoutOpts.CommitHash)
	// Try to resolve target commit
	commit, err = repo.CommitObject(targetHash)

	if err != nil && errors.Is(err, plumbing.ErrObjectNotFound) {
		// Need to fetch deeper
		depth := uint(0)
		maxDepth := gofn.Coalesce(checkoutOpts.MaxDepth, checkoutMaxDepthDefault)
		depthIncrement := max(20, maxDepth/10) //nolint:mnd

		env := []string{} // No use current process's env
		if checkoutOpts.sshCmd != "" {
			env = append(env, "GIT_SSH_COMMAND="+checkoutOpts.sshCmd)
		}

		for depth <= maxDepth {
			depth += depthIncrement
			fetchArgs := []string{"fetch", "origin", "--depth", strconv.FormatUint(uint64(checkoutOpts.Depth), 10)}
			if checkoutOpts.branch != "" {
				fetchArgs = append(fetchArgs, checkoutOpts.branch)
			}

			fetchCmd := exec.CommandContext(ctx, "git", fetchArgs...)
			fetchCmd.Dir = checkoutOpts.CheckoutPath
			fetchCmd.Env = env
			if out, fetchErr := fetchCmd.CombinedOutput(); fetchErr != nil {
				return nil, apperrors.New(fetchErr).WithExtraDetail("%s", string(out))
			}

			commit, err = repo.CommitObject(targetHash)
			if err == nil && commit != nil {
				break
			}
		}
	}

	if commit == nil {
		return nil, apperrors.Wrap(plumbing.ErrObjectNotFound)
	}

	// Checkout target commit
	checkoutCmd := exec.CommandContext(ctx, "git", "checkout", checkoutOpts.CommitHash) //nolint:gosec
	checkoutCmd.Dir = checkoutOpts.CheckoutPath
	checkoutCmd.Env = []string{} // No use current process's env
	if out, checkoutErr := checkoutCmd.CombinedOutput(); checkoutErr != nil {
		return nil, apperrors.New(checkoutErr).WithExtraDetail("%s", string(out))
	}

	return commit, nil
}

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

	env := []string{} // No use current process's env
	if checkoutOpts.sshCmd != "" {
		env = append(env, "GIT_SSH_COMMAND="+checkoutOpts.sshCmd)
	}

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = checkoutOpts.CheckoutPath
	cmd.Env = env
	if out, err := cmd.CombinedOutput(); err != nil {
		return apperrors.New(err).WithExtraDetail("%s", string(out))
	}
	return nil
}
