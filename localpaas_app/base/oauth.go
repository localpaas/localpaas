package base

type OAuthType string

const (
	OAuthTypeGithub = OAuthType("github")
	OAuthTypeGitlab = OAuthType("gitlab")
	OAuthTypeGoogle = OAuthType("google")
)

var (
	AllOAuthTypes = []OAuthType{OAuthTypeGithub, OAuthTypeGitlab, OAuthTypeGoogle}
)
