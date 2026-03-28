package taskcronjobexec

import (
	"context"
)

func (e *Executor) sysCleanupFiles(
	_ context.Context,
	data *sysCleanupTaskData,
) (err error) {
	defer func() {
		if err != nil {
			data.TaskOutput.FileCleanup.Error = err.Error()
		}
	}()

	// TODO: add implementation

	return nil
}
