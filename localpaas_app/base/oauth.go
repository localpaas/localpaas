package base

type OAuthType string

const (
	OAuthTypeGithub    OAuthType = "github"
	OAuthTypeGithubApp OAuthType = "github-app"
	OAuthTypeGitlab    OAuthType = "gitlab"
	OAuthTypeGitea     OAuthType = "gitea"
	OAuthTypeGoogle    OAuthType = "google"

	// Custom OAuth types
	OAuthTypeGitlabCustom OAuthType = "gitlab-custom"
)

var (
	AllOAuthTypes = []OAuthType{OAuthTypeGithub, OAuthTypeGithubApp, OAuthTypeGitlab, OAuthTypeGitea,
		OAuthTypeGoogle, OAuthTypeGitlabCustom}
)

func IsCustomOAuthType(typ OAuthType) bool {
	switch typ { //nolint:exhaustive
	case OAuthTypeGitlabCustom:
		return true
	default:
		return false
	}
}
