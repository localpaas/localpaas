package base

type ProjectStatus string

const (
	ProjectStatusActive   ProjectStatus = "active"
	ProjectStatusLocked   ProjectStatus = "locked"
	ProjectStatusDisabled ProjectStatus = "disabled"
	ProjectStatusDeleting ProjectStatus = "deleting"
)

var (
	AllProjectStatuses = []ProjectStatus{ProjectStatusActive, ProjectStatusLocked, ProjectStatusDisabled,
		ProjectStatusDeleting}
)
