package base

type GitSource string

const (
	GitSourceGithub    GitSource = "github"
	GitSourceGitlab    GitSource = "gitlab"
	GitSourceGitea     GitSource = "gitea"
	GitSourceBitbucket GitSource = "bitbucket"
)

var (
	AllGitSources = []GitSource{GitSourceGithub, GitSourceGitlab, GitSourceGitea, GitSourceBitbucket}
)
