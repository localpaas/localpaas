package base

import "time"

type CronJobType string

const (
	CronJobTypeContainerCommand CronJobType = "container-command"
)

var (
	AllCronJobTypes = []CronJobType{CronJobTypeContainerCommand}
)

const (
	// TODO: make this configurable
	CronJobNotificationTimeoutDefault = 1 * time.Minute
)
