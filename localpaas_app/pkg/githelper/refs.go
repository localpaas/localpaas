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
	if strings.HasPrefix(ref, "tags/") {
		ref = strings.TrimPrefix(ref, "tags/")
		return plumbing.NewTagReferenceName(ref)
	}

	// Heads ref
	if strings.HasPrefix(ref, "heads/") {
		ref = strings.TrimPrefix(ref, "heads/")
		return plumbing.NewBranchReferenceName(ref)
	}
	return plumbing.NewBranchReferenceName(ref)
}
