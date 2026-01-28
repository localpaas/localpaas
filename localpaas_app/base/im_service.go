package base

type IMServiceKind string

const (
	IMServiceKindSlack   IMServiceKind = "slack"
	IMServiceKindDiscord IMServiceKind = "discord"
)

var (
	AllIMServiceKinds = []IMServiceKind{IMServiceKindSlack, IMServiceKindDiscord}
)
