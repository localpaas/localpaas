package entity

type TaskCronJobNotificationArgs struct {
	App     ObjectID `json:"app"`
	CronJob ObjectID `json:"cronJob"`
}

type TaskCronJobNotificationOutput struct {
}

func (t *Task) ArgsAsCronJobNotification() (*TaskCronJobNotificationArgs, error) {
	return parseTaskArgsAs(t, func() *TaskCronJobNotificationArgs { return &TaskCronJobNotificationArgs{} })
}

func (t *Task) OutputAsCronJobNotification() (*TaskCronJobNotificationOutput, error) {
	return parseTaskOutputAs(t, func() *TaskCronJobNotificationOutput { return &TaskCronJobNotificationOutput{} })
}
