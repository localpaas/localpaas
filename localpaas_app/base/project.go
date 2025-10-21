package base

type ProjectStatus string

const (
	ProjectStatusActive   ProjectStatus = "active"
	ProjectStatusLocked   ProjectStatus = "locked"
	ProjectStatusDisabled ProjectStatus = "disabled"
)

var (
	AllProjectStatuses = []ProjectStatus{ProjectStatusActive, ProjectStatusLocked, ProjectStatusDisabled}
)
