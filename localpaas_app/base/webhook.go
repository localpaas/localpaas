package base

type WebhookKind string

const (
	WebhookKindGithub    WebhookKind = "github"
	WebhookKindGitlab    WebhookKind = "gitlab"
	WebhookKindGitea     WebhookKind = "gitea"
	WebhookKindBitbucket WebhookKind = "bitbucket"
	WebhookKindGogs      WebhookKind = "gogs"
)

var (
	AllWebhookKinds = []WebhookKind{WebhookKindGithub, WebhookKindGitlab, WebhookKindGitea,
		WebhookKindBitbucket, WebhookKindGogs}
)

const (
	DefaultWebhookSecretByteLen = 20 // string length should be double
)
