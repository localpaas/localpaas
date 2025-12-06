package base

type GitTokenType string

const (
	GitTokenTypeGithub       GitTokenType = "github"
	GitTokenTypeGitlab       GitTokenType = "gitlab"
	GitTokenTypeGitlabCustom GitTokenType = "gitlab-custom"
	GitTokenTypeGitea        GitTokenType = "gitea"
)

var (
	AllGitTokenTypes = []GitTokenType{GitTokenTypeGithub, GitTokenTypeGitlab, GitTokenTypeGitlabCustom,
		GitTokenTypeGitea}
)
