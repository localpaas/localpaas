package githelper

import (
	"strings"

	"github.com/go-git/go-git/v6/plumbing"
)

func NormalizeRepoRef(ref string) plumbing.ReferenceName {
	if ref == "" {
		return "HEAD"
	}
	if strings.HasPrefix(ref, "refs/") {
		return plumbing.ReferenceName(ref)
	}

	// Tags ref
	if after, ok := strings.CutPrefix(ref, "tags/"); ok {
		ref = after
		return plumbing.NewTagReferenceName(ref)
	}

	// Heads ref
	if after, ok := strings.CutPrefix(ref, "heads/"); ok {
		ref = after
		return plumbing.NewBranchReferenceName(ref)
	}
	return plumbing.NewBranchReferenceName(ref)
}
