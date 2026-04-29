package entity

type TaskAppVerUpdateArgs struct {
	Deployment ObjectID `json:"deployment"`
}

type TaskAppVerUpdateOutput struct {
}

func (t *Task) ArgsAsAppVerUpdate() (*TaskAppVerUpdateArgs, error) {
	return parseTaskArgsAs(t, func() *TaskAppVerUpdateArgs { return &TaskAppVerUpdateArgs{} })
}

func (t *Task) OutputAsAppVerUpdate() (*TaskAppVerUpdateOutput, error) {
	return parseTaskOutputAs(t, func() *TaskAppVerUpdateOutput { return &TaskAppVerUpdateOutput{} })
}
