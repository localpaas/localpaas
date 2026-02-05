package taskappdeploy

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
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/services/docker"
)

const (
	stepPreDeployCmd  = "pre-deploy-cmd-exec"
	stepPostDeployCmd = "post-deploy-cmd-exec"

	preDeploymentContainerFindRetryMax   = 1
	preDeploymentContainerFindRetryDelay = time.Second * 5

	postDeploymentContainerFindRetryMax   = 3
	postDeploymentContainerFindRetryDelay = time.Second * 10
)

func (e *Executor) deployStepExecCmd(
	ctx context.Context,
	data *taskData,
	preDeployment bool,
) (err error) {
	deployment := data.Deployment
	if preDeployment && (deployment.Settings.PreDeploymentCommand == nil ||
		*deployment.Settings.PreDeploymentCommand == "") {
		return nil
	} else {
		data.Step = stepPreDeployCmd
	}
	if !preDeployment && (deployment.Settings.PostDeploymentCommand == nil ||
		*deployment.Settings.PostDeploymentCommand == "") {
		return nil
	} else {
		data.Step = stepPostDeployCmd
	}

	e.addStepStartLog(ctx, data, fmt.Sprintf("Start executing %s-deployment command...",
		gofn.If(preDeployment, "pre", "post")))
	defer e.addStepEndLog(ctx, data, timeutil.NowUTC(), err)

	var maxRetry int
	var retryDelay time.Duration
	if preDeployment {
		maxRetry = preDeploymentContainerFindRetryMax
		retryDelay = preDeploymentContainerFindRetryDelay
	} else {
		maxRetry = postDeploymentContainerFindRetryMax
		retryDelay = postDeploymentContainerFindRetryDelay
	}

	contSum, _, err := e.dockerManager.ServiceContainerGetActive(ctx, data.App.ServiceID,
		maxRetry, retryDelay)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if contSum == nil {
		_ = data.LogStore.Add(ctx, realtimelog.NewWarnFrame(
			"No running container found, execution skipped", nil))
		return nil
	}

	var cmdStr string
	if preDeployment {
		cmdStr = *deployment.Settings.PreDeploymentCommand
	} else {
		cmdStr = *deployment.Settings.PostDeploymentCommand
	}

	execOptions := &container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          gofn.StringSplit(cmdStr, " ", "\""),
	}
	if deployment.Settings.WorkingDir != nil {
		execOptions.WorkingDir = *deployment.Settings.WorkingDir
	}

	execID, resp, err := e.dockerManager.ContainerExec(ctx, contSum.ID, execOptions)
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
		_ = data.LogStore.Add(ctx, realtimelog.NewWarnFrame(fmt.Sprintf(
			"Command execution failed with exit code: %v", execInfo.ExitCode), nil))
		return apperrors.Wrap(apperrors.ErrInfraActionFailed)
	}

	return nil
}
