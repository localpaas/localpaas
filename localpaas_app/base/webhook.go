package base

type WebhookKind string

const (
	WebhookKindGithub      WebhookKind = "github"
	WebhookKindGitlab      WebhookKind = "gitlab"
	WebhookKindGitea       WebhookKind = "gitea"
	WebhookKindBitbucket   WebhookKind = "bitbucket"
	WebhookKindGogs        WebhookKind = "gogs"
	WebhookKindAzureDevOps WebhookKind = "azuredevops"
)

var (
	AllWebhookKinds = []WebhookKind{WebhookKindGithub, WebhookKindGitlab, WebhookKindGitea,
		WebhookKindBitbucket, WebhookKindGogs, WebhookKindAzureDevOps}
)
