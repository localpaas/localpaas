package taskappdeploy

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/build"
	"github.com/gitsight/go-vcsurl"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/moby/go-archive"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/applog"
	"github.com/localpaas/localpaas/localpaas_app/pkg/batchrecvchan"
	"github.com/localpaas/localpaas/localpaas_app/pkg/fileutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/services/docker"
)

const (
	stepCodeCheckout = "code-checkout"
	stepImageBuild   = "image-build"
	stepServiceApply = "service-apply"
)

type repoDeployTaskData struct {
	*taskData
	CredSetting  *entity.Setting
	CheckoutPath string
	RepoURLInfo  *vcsurl.VCS
}

func (e *Executor) deployFromRepo(
	ctx context.Context,
	db database.Tx,
	taskData *taskData,
) error {
	data := &repoDeployTaskData{taskData: taskData}
	data.OnCommand(func(cmd base.TaskCommand, args ...any) { e.onRepoDeployCommand(ctx, data, cmd, args...) })

	// 0. Prepare
	err := e.repoDeployStepPrepare(ctx, db, data)
	if err != nil {
		return apperrors.Wrap(err)
	}
	defer e.repoDeployStepCleanup(ctx, data) //nolint:errcheck

	if data.IsCanceled() {
		return nil
	}

	// 1. Repo checkout
	err = e.repoDeployStepSourceCheckout(ctx, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	if data.IsCanceled() {
		return nil
	}

	// 2. Build image
	err = e.repoDeployStepImageBuild(ctx, db, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	if data.IsCanceled() {
		return nil
	}

	// 3. Pre-deployment command execution
	err = e.deployStepExecCmd(ctx, data.taskData, true)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// 4. Apply image to service
	err = e.repoDeployStepServiceApply(ctx, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// 5. Post-deployment command execution
	err = e.deployStepExecCmd(ctx, data.taskData, false)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (e *Executor) repoDeployStepSourceCheckout(
	ctx context.Context,
	data *repoDeployTaskData,
) (err error) {
	data.Step = stepCodeCheckout
	deployment := data.Deployment
	repoSource := deployment.Settings.RepoSource

	// NOTE: currently supports repo of git type only
	if repoSource.RepoType != base.RepoTypeGit {
		return apperrors.New(apperrors.ErrUnsupported).
			WithExtraDetail("Repo type %s is unsupported", repoSource.RepoType)
	}

	e.addStepStartLog(ctx, data.taskData, "Start cloning Git repository...")
	defer e.addStepEndLog(ctx, data.taskData, timeutil.NowUTC(), err)

	authMethod, err := e.calcGitAuthMethod(ctx, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	repo, err := git.PlainCloneContext(ctx, data.CheckoutPath, &git.CloneOptions{
		URL:               repoSource.RepoURL,
		ReferenceName:     plumbing.ReferenceName(repoSource.RepoRef),
		Auth:              authMethod,
		Depth:             1,
		SingleBranch:      true,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		ShallowSubmodules: true,
	})
	if err != nil {
		return apperrors.Wrap(err)
	}

	commitIter, err := repo.CommitObjects()
	if err != nil {
		return apperrors.Wrap(err)
	}
	commit, err := commitIter.Next()
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.DeploymentOutput.CommitHash = commit.Hash.String()
	data.DeploymentOutput.CommitMessage = commit.Message

	// Remove .git dir within the source dir
	ee := os.RemoveAll(filepath.Join(data.CheckoutPath, ".git"))
	if ee != nil { // Just log
		_ = data.LogStore.Add(ctx, applog.NewErrFrame("failed to remove .git folder",
			applog.TsNow))
	}

	return nil
}

func (e *Executor) repoDeployStepImageBuild(
	ctx context.Context,
	db database.Tx,
	data *repoDeployTaskData,
) (err error) {
	data.Step = stepImageBuild
	deployment := data.Deployment
	repoSource := deployment.Settings.RepoSource

	e.addStepStartLog(ctx, data.taskData, "Start building image...")
	defer e.addStepEndLog(ctx, data.taskData, timeutil.NowUTC(), err)

	// TODO: check dockerfile existence
	dockerfile := gofn.Coalesce(repoSource.DockerfilePath, "Dockerfile")

	imageTags := e.calcBuildImageTags(repoSource.ImageTags, data)
	data.DeploymentOutput.ImageTags = imageTags

	envVars, err := e.calcBuildEnvVars(ctx, db, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	authConfigs, err := e.calcBuildRegistryAuths(ctx, db, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Create tar archive for the source code
	tar, err := archive.TarWithOptions(data.CheckoutPath, &archive.TarOptions{})
	if err != nil {
		return apperrors.Wrap(err)
	}
	defer tar.Close()

	// Build the image
	resp, err := e.dockerManager.ImageBuild(ctx, tar, func(opts *build.ImageBuildOptions) {
		opts.Version = build.BuilderV1
		opts.BuildID = data.Task.ID
		opts.Dockerfile = dockerfile
		opts.Tags = imageTags
		opts.BuildArgs = envVars
		opts.AuthConfigs = authConfigs
	})
	if err != nil {
		return apperrors.Wrap(err)
	}

	logsChan, _ := docker.StartScanningJSONMsg(ctx, resp.Body, batchrecvchan.Options{})
	for msgs := range logsChan {
		for _, msg := range msgs {
			if msg.Error != nil {
				err = errors.Join(err, msg.Error)
				_ = data.LogStore.Add(ctx, applog.NewErrFrame(msg.String(), applog.TsNow))
			} else {
				_ = data.LogStore.Add(ctx, applog.NewOutFrame(msg.String(), applog.TsNow))
			}
		}
	}
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (e *Executor) repoDeployStepServiceApply(
	ctx context.Context,
	data *repoDeployTaskData,
) (err error) {
	data.Step = stepServiceApply
	deployment := data.Deployment

	e.addStepStartLog(ctx, data.taskData, "Applying changes to service...")
	defer e.addStepEndLog(ctx, data.taskData, timeutil.NowUTC(), err)

	service, err := e.dockerManager.ServiceInspect(ctx, data.App.ServiceID)
	if err != nil {
		return apperrors.Wrap(err)
	}

	spec := &service.Spec
	contSpec := spec.TaskTemplate.ContainerSpec
	contSpec.Image = data.DeploymentOutput.ImageTags[0]
	if deployment.Settings.WorkingDir != nil {
		contSpec.Dir = *deployment.Settings.WorkingDir
	}
	if deployment.Settings.Command != nil {
		docker.ApplyServiceCommand(contSpec, *deployment.Settings.Command)
	}

	_, err = e.dockerManager.ServiceUpdate(ctx, data.App.ServiceID, &service.Version, spec)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (e *Executor) repoDeployStepPrepare(
	ctx context.Context,
	db database.Tx,
	data *repoDeployTaskData,
) (err error) {
	deployment := data.Deployment
	repoSource := deployment.Settings.RepoSource

	// Loads repo credentials (github app, git token, ssh key) if configured
	settingID := repoSource.Credentials.ID
	if settingID != "" {
		setting, err := e.settingRepo.GetByID(ctx, db, "", settingID, true)
		if err != nil {
			return apperrors.Wrap(err)
		}
		data.CredSetting = setting
	}

	// Creates checkout dir
	data.CheckoutPath, err = fileutil.CreateTempDir("", "*", 0)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Parse repo URL
	repoURLInfo, err := vcsurl.Parse(repoSource.RepoURL)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.RepoURLInfo = repoURLInfo

	return nil
}

//nolint:unparam
func (e *Executor) repoDeployStepCleanup(
	_ context.Context,
	data *repoDeployTaskData,
) (err error) {
	if data.CheckoutPath != "" {
		_ = os.RemoveAll(data.CheckoutPath)
	}

	return nil
}

func (e *Executor) onRepoDeployCommand(
	ctx context.Context,
	taskData *repoDeployTaskData,
	cmd base.TaskCommand,
	_ ...any,
) {
	if cmd == base.TaskCommandCancel && taskData.Step == stepImageBuild {
		err := e.dockerManager.ImageBuildCancel(ctx, taskData.Task.ID)
		if err != nil {
			_ = taskData.LogStore.Add(ctx, applog.NewErrFrame("failed to cancel image build: "+
				err.Error(), applog.TsNow))
		}
	}
}
