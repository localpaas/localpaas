package taskappdeploy

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/registry"
	"github.com/go-git/go-git/v6/plumbing/transport"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
	"github.com/go-git/go-git/v6/plumbing/transport/ssh"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
	"github.com/localpaas/localpaas/services/git/github"
)

func (e *Executor) calcGitAuthMethod(
	ctx context.Context,
	data *repoDeployTaskData,
) (auth transport.AuthMethod, err error) {
	if data.CredSetting == nil {
		return auth, nil
	}
	switch data.CredSetting.Type { //nolint:exhaustive
	case base.SettingTypeGithubApp:
		client, err := github.NewFromSetting(data.CredSetting)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		token, err := client.CreateAppToken(ctx)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		auth = &http.BasicAuth{
			Username: "default", // this can be anything except an empty string
			Password: token,
		}

	case base.SettingTypeAccessToken:
		token, err := data.CredSetting.MustAsAccessToken().Token.GetPlain()
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
	return auth, nil
}

func (e *Executor) calcBuildImageTags(
	imageTags []string,
	data *repoDeployTaskData,
) ([]string, error) {
	if len(imageTags) > 0 {
		return imageTags, nil
	}

	imageName := data.Deployment.Settings.RepoSource.ImageName
	if imageName == "" || imageName == "auto" {
		imageName = data.App.GetAutoImageName()
	}

	commitHashPortion := data.DeploymentOutput.CommitHash[:7]
	tagCurrent := fmt.Sprintf("%s:%s", imageName, commitHashPortion)

	// If `pushToRegistry` is set in the settings, need to prepend the registry domain and
	// username to the tags.
	// E.g. `app_name:latest` will likely become `docker.io/username/app_name:latest`
	repoSource := data.Deployment.Settings.RepoSource
	if repoSource.PushToRegistry.ID != "" {
		regAuthSetting := data.RefSettingMap[repoSource.PushToRegistry.ID]
		if regAuthSetting == nil {
			return nil, apperrors.NewMissing("Registry auth to push image")
		}
		regAuth := regAuthSetting.MustAsRegistryAuth()
		tagCurrentWithReg := regAuth.Address + "/" + regAuth.Username + "/" + tagCurrent
		imageTags = append(imageTags, tagCurrentWithReg)
	}

	imageTags = append(imageTags, tagCurrent)
	return imageTags, nil
}

func (e *Executor) calcBuildEnvVars(
	ctx context.Context,
	db database.Tx,
	data *repoDeployTaskData,
) (map[string]*string, error) {
	envVars, err := e.envVarService.BuildAppEnvVars(ctx, db, data.App, true)
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
	app := data.App

	settings, _, err := e.settingRepo.ListByProject(ctx, db, app.ProjectID, nil,
		bunex.SelectWhere("setting.type = ?", base.SettingTypeRegistryAuth),
		bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
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
