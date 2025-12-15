package base

type TaskType string

const (
	TaskTypeSchedule TaskType = "task:schedule"
	TaskTypeTest     TaskType = "task:test"
	TaskTypeGitClone TaskType = "task:git-clone"
)

var (
	AllTaskTypes = []TaskType{TaskTypeSchedule, TaskTypeTest, TaskTypeGitClone}
)

type TaskStatus string

const (
	TaskStatusNotStarted TaskStatus = "not-started"
	TaskStatusInProgress TaskStatus = "in-progress"
	TaskStatusCanceled   TaskStatus = "canceled"
	TaskStatusDone       TaskStatus = "done"
	TaskStatusFailed     TaskStatus = "failed"
)

var (
	AllTaskStatuses = []TaskStatus{TaskStatusNotStarted, TaskStatusInProgress, TaskStatusCanceled,
		TaskStatusDone, TaskStatusFailed}
	AllTaskSettableStatuses = []TaskStatus{TaskStatusCanceled}
)

type TaskPriority string

const (
	TaskPriorityLow      TaskPriority = "low"
	TaskPriorityDefault  TaskPriority = "default"
	TaskPriorityCritical TaskPriority = "critical"
)

var (
	AllTaskPriorities = []TaskPriority{TaskPriorityLow, TaskPriorityDefault, TaskPriorityCritical}
)
