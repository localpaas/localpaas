package taskappdeploy

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/registry"
	"github.com/gitsight/go-vcsurl"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/transport"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
	"github.com/go-git/go-git/v6/plumbing/transport/ssh"
	"github.com/moby/go-archive"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/batchrecvchan"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/fileutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
	"github.com/localpaas/localpaas/services/docker"
)

const (
	maxImageNameLen = 200
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
	_ context.Context,
	data *repoDeployTaskData,
) error {
	deployment := data.Deployment
	repoSource := deployment.Settings.RepoSource

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
	err = os.RemoveAll(filepath.Join(data.CheckoutPath, ".git"))
	if err != nil {
		e.logger.Warn("failed to remove .git folder") // Just log
	}

	return nil
}

func (e *Executor) repoDeployStepBuildImage(
	ctx context.Context,
	db database.Tx,
	data *repoDeployTaskData,
) error {
	deployment := data.Deployment
	repoSource := deployment.Settings.RepoSource

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
		// print(">>>>>>>>>> ", reflectutil.UnsafeBytesToStr(gofn.Must(json.Marshal(msg))))
		if msgs[0].Error != nil {
			err = errors.Join(err, msgs[0].Error)
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
) error {
	deployment := data.Deployment

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

func (e *Executor) calcGitRefName(ref string) plumbing.ReferenceName {
	if ref == "" {
		return "HEAD"
	}

	// Tags ref
	if strings.HasPrefix(ref, "tags/") {
		ref = strings.TrimPrefix(ref, "tags/")
		return plumbing.NewTagReferenceName(ref)
	}
	if strings.HasPrefix(ref, "refs/tags/") {
		ref = strings.TrimPrefix(ref, "refs/tags/")
		return plumbing.NewTagReferenceName(ref)
	}

	// Heads ref
	if strings.HasPrefix(ref, "heads/") {
		ref = strings.TrimPrefix(ref, "heads/")
		return plumbing.NewBranchReferenceName(ref)
	}
	if strings.HasPrefix(ref, "refs/heads/") {
		ref = strings.TrimPrefix(ref, "refs/heads/")
		return plumbing.NewBranchReferenceName(ref)
	}
	return plumbing.NewBranchReferenceName(ref)
}

func (e *Executor) calcGitAuthMethod(
	data *repoDeployTaskData,
) (auth transport.AuthMethod, err error) {
	if data.CredSetting != nil {
		switch data.CredSetting.Type { //nolint:exhaustive
		case base.SettingTypeGithubApp:
			// TODO: add implementation
			break

		case base.SettingTypeGitToken:
			token, err := data.CredSetting.MustAsGitToken().Token.GetPlain()
			if err != nil {
				return nil, apperrors.Wrap(err)
			}
			auth = &http.BasicAuth{
				Username: "default", // this can be anything except an empty string
				Password: token,
			}

		case base.SettingTypeSSHKey:
			sshKey := data.CredSetting.MustAsSSHKey()
			privateKey, err := sshKey.PrivateKey.GetPlain()
			if err != nil {
				return nil, apperrors.Wrap(err)
			}
			passphrase, err := sshKey.Passphrase.GetPlain()
			if err != nil {
				return nil, apperrors.Wrap(err)
			}
			gitUser := gofn.Coalesce(data.RepoURLInfo.Username, "git")
			auth, err = ssh.NewPublicKeys(gitUser, reflectutil.UnsafeStrToBytes(privateKey), passphrase)
			if err != nil {
				return nil, apperrors.Wrap(err)
			}
		}
	}
	return auth, nil
}

func (e *Executor) calcBuildImageTags(
	imageTags []string,
	data *repoDeployTaskData,
) []string {
	if len(imageTags) > 0 {
		return imageTags
	}

	appKey := data.Deployment.App.Key
	if len(appKey) > maxImageNameLen {
		appKey = appKey[:maxImageNameLen]
	}

	commitHashPortion := data.DeploymentOutput.CommitHash[:7]
	imageTags = append(imageTags, fmt.Sprintf("%s:latest", appKey),
		fmt.Sprintf("%s:%s", appKey, commitHashPortion))

	return imageTags
}

func (e *Executor) calcBuildEnvVars(
	ctx context.Context,
	db database.Tx,
	data *repoDeployTaskData,
) (map[string]*string, error) {
	envVars, err := e.envVarService.BuildAppEnv(ctx, db, data.Deployment.App, true)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	result := make(map[string]*string, len(envVars))
	for _, envVar := range envVars {
		result[envVar.Key] = &envVar.Value
	}

	return result, nil
}

func (e *Executor) calcBuildRegistryAuths(
	ctx context.Context,
	db database.Tx,
	data *repoDeployTaskData,
) (map[string]registry.AuthConfig, error) {
	deployment := data.Deployment

	settings, _, err := e.settingRepo.List(ctx, db, nil,
		bunex.SelectWhere("setting.type = ?", base.SettingTypeRegistryAuth),
		bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
		bunex.SelectJoin("LEFT JOIN project_shared_settings pss ON pss.setting_id = setting.id"),
		bunex.SelectWhereGroup(
			bunex.SelectWhere("setting.object_id = ?", deployment.App.ProjectID),
			bunex.SelectWhereOr("pss.project_id = ?", deployment.App.ProjectID),
		),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	result := make(map[string]registry.AuthConfig, len(settings))
	for _, setting := range settings {
		regAuth, err := setting.AsRegistryAuth()
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		password, err := regAuth.Password.GetPlain()
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		result[regAuth.Address] = registry.AuthConfig{
			Username:      regAuth.Username,
			Password:      password,
			ServerAddress: regAuth.Address,
		}
	}

	return result, nil
}
