package entity

type TaskAppNotificationArgs struct {
	App        ObjectID  `json:"app"`
	Deployment *ObjectID `json:"deployment"`
}

type TaskAppNotificationOutput struct {
}

func (t *Task) ArgsAsAppNotification() (*TaskAppNotificationArgs, error) {
	return parseTaskArgsAs(t, func() *TaskAppNotificationArgs { return &TaskAppNotificationArgs{} })
}

func (t *Task) OutputAsAppNotification() (*TaskAppNotificationOutput, error) {
	return parseTaskOutputAs(t, func() *TaskAppNotificationOutput { return &TaskAppNotificationOutput{} })
}
