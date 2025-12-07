package base

type OAuthKind string

const (
	OAuthKindGithub    OAuthKind = "github"
	OAuthKindGithubApp OAuthKind = "github-app"
	OAuthKindGitlab    OAuthKind = "gitlab"
	OAuthKindGitea     OAuthKind = "gitea"
	OAuthKindGoogle    OAuthKind = "google"

	// Custom OAuth types
	OAuthKindGitlabCustom OAuthKind = "gitlab-custom"
)

var (
	AllOAuthKinds = []OAuthKind{OAuthKindGithub, OAuthKindGithubApp, OAuthKindGitlab, OAuthKindGitea,
		OAuthKindGoogle, OAuthKindGitlabCustom}
)

func IsCustomOAuthKind(kind OAuthKind) bool {
	switch kind { //nolint:exhaustive
	case OAuthKindGitlabCustom:
		return true
	default:
		return false
	}
}
