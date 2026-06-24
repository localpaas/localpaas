package githelper

import (
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
)

const (
	refHeadsPrefix         = "refs/heads/"
	refTagsPrefix          = "refs/tags/"
	refPullPrefix          = "refs/pull/"
	refMergeRequestsPrefix = "refs/merge-requests/"
)

type RefType string

const (
	RefBranch RefType = "branch"
	RefTag    RefType = "tag"
	RefPull   RefType = "pull"
)

func (rt RefType) IsBranch() bool {
	return rt == RefBranch
}

func (rt RefType) IsTag() bool {
	return rt == RefTag
}

func (rt RefType) IsPull() bool {
	return rt == RefPull
}

func (rt RefType) CanCheckout() bool {
	return rt == RefBranch || rt == RefTag || rt == RefPull
}

func NormalizeRepoRef(ref string) plumbing.ReferenceName {
	if ref == "" || ref == "HEAD" { //nolint:goconst
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

	// Pull ref (github, gitea)
	if after, ok := strings.CutPrefix(ref, "pull/"); ok {
		ref = after
		ref, _ = strings.CutSuffix(ref, "/head")
		return plumbing.ReferenceName(refPullPrefix + ref + "/head")
	}

	// Merge request ref (gitlab)
	if after, ok := strings.CutPrefix(ref, "merge-requests/"); ok {
		ref = after
		ref, _ = strings.CutSuffix(ref, "/head")
		return plumbing.ReferenceName(refMergeRequestsPrefix + ref + "/head")
	}

	// Branch
	return plumbing.NewBranchReferenceName(ref)
}

func GetRefType(ref string) RefType {
	if strings.HasPrefix(ref, refHeadsPrefix) {
		return RefBranch
	}
	if strings.HasPrefix(ref, refTagsPrefix) {
		return RefTag
	}
	if strings.HasPrefix(ref, refPullPrefix) || strings.HasPrefix(ref, refMergeRequestsPrefix) {
		return RefPull
	}
	return ""
}

func GetRefShort(ref string) (RefType, string) {
	refType := GetRefType(ref)
	if refType == RefBranch || refType == RefTag {
		return refType, plumbing.ReferenceName(ref).Short()
	}
	return refType, ref
}
