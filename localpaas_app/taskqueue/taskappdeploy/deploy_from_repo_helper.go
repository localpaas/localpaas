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
)

const (
	maxImageNameLen = 200
)

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
	app := deployment.App

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
