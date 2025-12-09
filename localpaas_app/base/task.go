package base

type TaskType string

type TaskStatus string

const (
	TaskStatusNotStarted TaskStatus = "not-started"
	TaskStatusInProgress TaskStatus = "in-progress"
	TaskStatusCanceled   TaskStatus = "canceled"
	TaskStatusDone       TaskStatus = "done"
)

var (
	AllTaskStatuses = []TaskStatus{TaskStatusNotStarted, TaskStatusInProgress, TaskStatusCanceled,
		TaskStatusDone}
)
