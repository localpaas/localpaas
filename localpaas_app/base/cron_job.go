package base

import "time"

type CronJobType string

const (
	CronJobTypeContainerCommand CronJobType = "container-command"
	CronJobTypeSystemCleanup    CronJobType = "system-cleanup"
)

var (
	AllCronJobTypes = []CronJobType{CronJobTypeContainerCommand, CronJobTypeSystemCleanup}
)

const (
	// TODO: make this configurable
	CronTaskNotificationTimeoutDefault = 1 * time.Minute
)
