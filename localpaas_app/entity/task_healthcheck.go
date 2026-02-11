package entity

type TaskHealthcheckArgs struct {
}

type TaskHealthcheckOutput struct {
	REST *TaskHealthcheckOutputREST `json:"rest,omitempty"`
	GRPC *TaskHealthcheckOutputGRPC `json:"grpc,omitempty"`
}

type TaskHealthcheckOutputREST struct {
	ReturnCode int    `json:"returnCode,omitempty"`
	ReturnText string `json:"returnText,omitempty"`
	ReturnJSON string `json:"returnJSON,omitempty"`
}

type TaskHealthcheckOutputGRPC struct {
	// TODO: implement this
}

func (t *Task) ArgsAsHealthcheck() (*TaskHealthcheckArgs, error) {
	return parseTaskArgsAs(t, func() *TaskHealthcheckArgs { return &TaskHealthcheckArgs{} })
}

func (t *Task) OutputAsHealthcheck() (*TaskHealthcheckOutput, error) {
	return parseTaskOutputAs(t, func() *TaskHealthcheckOutput { return &TaskHealthcheckOutput{} })
}
