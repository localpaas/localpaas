package base

import "github.com/tiendc/gofn"

type OAuthType string

const (
	OAuthTypeGithub = OAuthType("github")
	OAuthTypeGitlab = OAuthType("gitlab")
	OAuthTypeGoogle = OAuthType("google")

	// Custom OAuth types
	OAuthTypeGitlabCustom = OAuthType("gitlab-custom")
)

var (
	AllOAuthTypes = []OAuthType{OAuthTypeGithub, OAuthTypeGitlab, OAuthTypeGoogle,
		OAuthTypeGitlabCustom}
)

func IsCustomOAuthType(typ OAuthType) bool {
	switch typ { //nolint:exhaustive
	case OAuthTypeGitlabCustom:
		return true
	default:
		return false
	}
}

func IsValidOAuthType(typ OAuthType) bool {
	return gofn.Contain(AllOAuthTypes, typ)
}
