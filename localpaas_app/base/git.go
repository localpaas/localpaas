package base

type GitSource string

const (
	GitSourceGithub       GitSource = "github"
	GitSourceGitlab       GitSource = "gitlab"
	GitSourceGitlabCustom GitSource = "gitlab-custom"
	GitSourceGitea        GitSource = "gitea"
)

var (
	AllGitSources = []GitSource{GitSourceGithub, GitSourceGitlab, GitSourceGitlabCustom,
		GitSourceGitea}
)
