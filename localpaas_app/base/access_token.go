package base

type TokenKind string

const (
	TokenKindGithub    TokenKind = "github"
	TokenKindGitlab    TokenKind = "gitlab"
	TokenKindGitea     TokenKind = "gitea"
	TokenKindBitbucket TokenKind = "bitbucket"
	TokenKindGogs      TokenKind = "gogs"
)

var (
	AllGitTokenKinds = []TokenKind{TokenKindGithub, TokenKindGitlab, TokenKindGitea, TokenKindBitbucket,
		TokenKindGogs}

	AllTokenKinds = append([]TokenKind{}, AllGitTokenKinds...)
)
