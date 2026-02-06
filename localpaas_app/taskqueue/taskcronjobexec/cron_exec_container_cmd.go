package taskcronjobexec

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/pkg/batchrecvchan"
	"github.com/localpaas/localpaas/localpaas_app/pkg/realtimelog"
	"github.com/localpaas/localpaas/services/docker"
)

const (
	retryMax   = 3
	retryDelay = time.Second * 5
)

func (e *Executor) cronExecContainerCmd(
	ctx context.Context,
	data *taskData,
) (err error) {
	command := data.CronJob.Command
	if command == nil || command.Command == "" { // can't continue if this happens
		data.NonRetryable = true
		data.Logs = append(data.Logs, realtimelog.NewErrFrame(
			"execution command is empty, execution aborted", nil))
		return apperrors.New(apperrors.ErrInternalServer).WithMsgLog("cron job command is empty")
	}

	contSum, _, err := e.dockerManager.ServiceContainerGetActive(ctx, data.App.ServiceID,
		retryMax, retryDelay)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if contSum == nil {
		data.Logs = append(data.Logs, realtimelog.NewWarnFrame(
			"No running container found, execution skipped", nil))
		return nil
	}

	execOptions := &container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          gofn.StringSplit(command.Command, " ", "\""),
		WorkingDir:   command.WorkingDir,
	}

	execID, resp, err := e.dockerManager.ContainerExec(ctx, contSum.ID, execOptions)
	if err != nil {
		return apperrors.Wrap(err)
	}

	logChan, _ := docker.StartScanningLog(ctx, io.NopCloser(resp.Reader), batchrecvchan.Options{})
	defer resp.Close()

	for msgs := range logChan {
		data.Logs = append(data.Logs, msgs...)
	}

	// Get exit code
	execInfo, err := e.dockerManager.ContainerExecInspect(ctx, execID)
	if err != nil {
		return apperrors.Wrap(err)
	}

	if execInfo.ExitCode != 0 {
		data.Logs = append(data.Logs, realtimelog.NewWarnFrame(fmt.Sprintf(
			"Command execution failed with exit code: %v", execInfo.ExitCode), nil))
		return apperrors.Wrap(apperrors.ErrInfraActionFailed)
	}

	return nil
}
