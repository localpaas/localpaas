package base

type SystemSettingKind string

const (
	SystemSettingKindDBCleanup SystemSettingKind = "db-cleanup"
)

var (
	AllSystemSettingKinds = []SystemSettingKind{SystemSettingKindDBCleanup}
)
