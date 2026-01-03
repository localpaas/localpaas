package taskappdeploy

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/pkg/batchrecvchan"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/services/docker"
)

const (
	stepPreDeployCmd  = "pre-deploy-cmd-exec"
	stepPostDeployCmd = "post-deploy-cmd-exec"
)

func (e *Executor) deployStepExecCmd(
	ctx context.Context,
	data *taskData,
	preDeployment bool,
) (err error) {
	deployment := data.Deployment
	if preDeployment && (deployment.Settings.PreDeployment == nil ||
		deployment.Settings.PreDeployment.Cmd == "") {
		return nil
	}
	if !preDeployment && (deployment.Settings.PostDeployment == nil ||
		deployment.Settings.PostDeployment.Cmd == "") {
		return nil
	}

	e.addStepStartLog(ctx, data, fmt.Sprintf("Start executing %s-deployment command...",
		gofn.If(preDeployment, "pre", "post")))
	defer e.addStepEndLog(ctx, data, timeutil.NowUTC(), err)

	containerSum, err := e.findContainerForCmdExec(ctx, data)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if containerSum == nil {
		return nil // TODO: handle this
	}

	var cmdStr string
	if preDeployment {
		cmdStr = deployment.Settings.PreDeployment.Cmd
	} else {
		cmdStr = deployment.Settings.PostDeployment.Cmd
	}
	cmd := gofn.StringSplit(cmdStr, " ", "\"")

	execID, resp, err := e.dockerManager.ContainerExec(ctx, containerSum.ID, &container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	})
	if err != nil {
		return apperrors.Wrap(err)
	}

	logChan, _ := docker.StartScanningLog(ctx, io.NopCloser(resp.Reader), batchrecvchan.Options{})
	defer resp.Close()

	for msgs := range logChan {
		_ = data.LogStore.Add(ctx, msgs...)
	}

	// Get exit code
	execInfo, err := e.dockerManager.ContainerExecInspect(ctx, execID)
	if err != nil {
		return apperrors.Wrap(err)
	}

	if execInfo.ExitCode != 0 {
		// TODO: handle this
	}

	return nil
}

func (e *Executor) findContainerForCmdExec(
	ctx context.Context,
	data *taskData,
) (*container.Summary, error) {
	app := data.Deployment.App

	containers, err := e.dockerManager.ServiceContainerList(ctx, app.ServiceID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	for i := range containers {
		c := &containers[i]
		if c.State == container.StateRunning {
			return c, nil
		}
	}

	return nil, nil
}
