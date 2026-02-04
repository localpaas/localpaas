package base

type OAuthKind string

const (
	OAuthKindGithub    OAuthKind = "github"
	OAuthKindGithubApp OAuthKind = "github-app"
	OAuthKindGitlab    OAuthKind = "gitlab"
	OAuthKindGitea     OAuthKind = "gitea"
	OAuthKindGoogle    OAuthKind = "google"
)

var (
	AllOAuthKinds = []OAuthKind{OAuthKindGithub, OAuthKindGithubApp, OAuthKindGitlab, OAuthKindGitea,
		OAuthKindGoogle}
)
