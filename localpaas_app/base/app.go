package base

type AppStatus string

const (
	AppStatusActive   AppStatus = "active"
	AppStatusLocked   AppStatus = "locked"
	AppStatusDisabled AppStatus = "disabled"
	AppStatusDeleting AppStatus = "deleting"
)

var (
	AllAppStatuses = []AppStatus{AppStatusActive, AppStatusLocked, AppStatusDisabled, AppStatusDeleting}
)
