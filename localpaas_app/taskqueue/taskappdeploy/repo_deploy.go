package taskappdeploy

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/build"
	"github.com/gitsight/go-vcsurl"
	"github.com/go-git/go-git/v6"
	"github.com/moby/go-archive"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/batchrecvchan"
	"github.com/localpaas/localpaas/localpaas_app/pkg/fileutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/realtimelog"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/services/docker"
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

	// 0. Prepare
	err := e.repoDeployStepPrepare(ctx, db, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// 1. Repo checkout
	err = e.repoDeployStepGitCheckout(ctx, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// 1.1. Check if deployment is canceled by user while we are processing it
	isCanceled, err := e.checkDeploymentCanceled(ctx, data.taskData)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if isCanceled {
		return nil
	}

	// 2. Build image
	err = e.repoDeployStepBuildImage(ctx, db, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// 2.1. Check if deployment is canceled by user while we are processing it
	isCanceled, err = e.checkDeploymentCanceled(ctx, data.taskData)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if isCanceled {
		return nil
	}

	// 3. Apply image to service
	err = e.repoDeployStepUpdateService(ctx, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (e *Executor) repoDeployStepGitCheckout(
	ctx context.Context,
	data *repoDeployTaskData,
) (err error) {
	deployment := data.Deployment
	repoSource := deployment.Settings.RepoSource

	e.addStepStartLogs(ctx, data.taskData, "Start cloning Git repository...")
	defer e.addStepEndLogs(ctx, data.taskData, timeutil.NowUTC(), err)

	// if args.Timeout > 0 {
	//	var cancel context.CancelFunc
	//	ctx, cancel = context.WithTimeout(ctx, args.Timeout)
	//	defer cancel()
	// }

	authMethod, err := e.calcGitAuthMethod(data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	repo, err := git.PlainClone(data.CheckoutPath, &git.CloneOptions{
		URL:               repoSource.RepoURL,
		ReferenceName:     e.calcGitRefName(repoSource.RepoRef),
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
		_ = data.LogStore.Add(ctx, realtimelog.NewErrFrame("failed to remove .git folder", nil))
	}

	return nil
}

func (e *Executor) repoDeployStepBuildImage(
	ctx context.Context,
	db database.Tx,
	data *repoDeployTaskData,
) (err error) {
	deployment := data.Deployment
	repoSource := deployment.Settings.RepoSource

	e.addStepStartLogs(ctx, data.taskData, "Start building image...")
	defer e.addStepEndLogs(ctx, data.taskData, timeutil.NowUTC(), err)

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
		// opts.BuildID = deployment.ID
		opts.Version = build.BuilderV1

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
				_ = data.LogStore.Add(ctx, realtimelog.NewErrFrame(msg.String(), nil))
			} else {
				_ = data.LogStore.Add(ctx, realtimelog.NewOutFrame(msg.String(), nil))
			}
		}
	}
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (e *Executor) repoDeployStepUpdateService(
	ctx context.Context,
	data *repoDeployTaskData,
) (err error) {
	deployment := data.Deployment

	e.addStepStartLogs(ctx, data.taskData, "Applying changes to service...")
	defer e.addStepEndLogs(ctx, data.taskData, timeutil.NowUTC(), err)

	service, err := e.dockerManager.ServiceInspect(ctx, deployment.App.ServiceID)
	if err != nil {
		return apperrors.Wrap(err)
	}

	spec := &service.Spec
	spec.TaskTemplate.ContainerSpec.Image = data.DeploymentOutput.ImageTags[0]

	_, err = e.dockerManager.ServiceUpdate(ctx, deployment.App.ServiceID, &service.Version, spec)
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
