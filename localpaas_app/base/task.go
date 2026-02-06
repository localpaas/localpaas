package base

type TaskType string

const (
	TaskTypeTest                TaskType = "task:test"
	TaskTypeAppDeploy           TaskType = "task:app-deploy"
	TaskTypeAppNotification     TaskType = "task:app-notification"
	TaskTypeCronJobExec         TaskType = "task:cron-job-exec"
	TaskTypeCronJobNotification TaskType = "task:cron-job-notification"
)

var (
	AllTaskTypes = []TaskType{TaskTypeTest, TaskTypeAppDeploy, TaskTypeAppNotification,
		TaskTypeCronJobExec, TaskTypeCronJobNotification}
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

	//nolint:mnd
	mapPriorityValues = map[TaskPriority]int{
		TaskPriorityLow:      3,
		TaskPriorityDefault:  6,
		TaskPriorityCritical: 10,
	}
)

func (p TaskPriority) Cmp(priority TaskPriority) int {
	if priority == "" {
		priority = TaskPriorityDefault
	}
	if p == priority {
		return 0
	}
	return mapPriorityValues[p] - mapPriorityValues[priority]
}

type TaskCommand string

const (
	TaskCommandCancel TaskCommand = "cancel"
)
