package base

type WebhookKind string

const (
	WebhookKindGithub    WebhookKind = "github"
	WebhookKindGitlab    WebhookKind = "gitlab"
	WebhookKindGitea     WebhookKind = "gitea"
	WebhookKindBitbucket WebhookKind = "bitbucket"
)

var (
	AllWebhookKinds = []WebhookKind{WebhookKindGithub, WebhookKindGitlab, WebhookKindGitea,
		WebhookKindBitbucket}
)
