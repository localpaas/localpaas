package appdeploymentserviceimpl

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/moby/moby/api/types/registry"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/githelper"
	"github.com/localpaas/localpaas/localpaas_app/pkg/redact"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/tasklog"
	"github.com/localpaas/localpaas/services/git/github"
)

func (s *service) calcGitAuthMethod(
	ctx context.Context,
	data *repoDeploymentData,
) (auth transport.AuthMethod, err error) {
	credSetting := data.CredSetting
	if credSetting == nil {
		return auth, nil
	}
	switch credSetting.Type { //nolint:exhaustive
	case base.SettingTypeGithubApp:
		client, err := github.NewFromSetting(credSetting)
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
		token, err := credSetting.MustAsAccessToken().Token.GetPlain()
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		auth = &http.BasicAuth{
			Username: "default", // this can be anything except an empty string
			Password: token,
		}

	case base.SettingTypeSSHKey:
		sshKey := credSetting.MustAsSSHKey()
		privateKey, err := sshKey.PrivateKey.GetPlain()
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		passphrase, err := sshKey.Passphrase.GetPlain()
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		authRaw, err := ssh.NewPublicKeys("git", reflectutil.UnsafeStrToBytes(privateKey), passphrase)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		auth = &githelper.AuthSSHKey{
			PublicKeys: authRaw,
			PEMBytes:   reflectutil.UnsafeStrToBytes(privateKey),
		}
	}
	return auth, nil
}

func (s *service) calcBuildImageTags(
	imageTags []string,
	data *repoDeploymentData,
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
		regAuthSetting := data.RefObjects.RefSettings[repoSource.PushToRegistry.ID]
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

func (s *service) calcBuildEnvVars(
	ctx context.Context,
	db database.Tx,
	data *repoDeploymentData,
) (map[string]*string, error) {
	envVars, refSecrets, err := s.envVarService.BuildAppEnvVars(ctx, db, data.App, true)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	data.SecretsToRedact = refSecrets

	if data.LogStore != nil && len(refSecrets) > 0 {
		secrets := make([]string, 0, len(refSecrets))
		for _, secret := range refSecrets {
			secrets = append(secrets, secret.Value.MustGetPlain())
		}
		data.LogStore.SetRedactor(redact.New(secrets))
	}

	result := make(map[string]*string, len(envVars))
	for _, envVar := range envVars {
		result[envVar.Key] = &envVar.Value
	}

	return result, nil
}

func (s *service) calcBuildRegistryAuths(
	ctx context.Context,
	db database.Tx,
	data *repoDeploymentData,
) (map[string]registry.AuthConfig, error) {
	app := data.App

	settings, _, err := s.settingRepo.List(ctx, db, base.NewObjectScopeProject(app.ProjectID), nil,
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

func (s *service) loadImageBuildSettings(
	ctx context.Context,
	db database.IDB,
	data *repoDeploymentData,
) error {
	app := data.App
	setting, err := s.settingRepo.GetSingle(ctx, db, base.NewObjectScopeProject(app.ProjectID),
		base.SettingTypeImageBuildSettings, true)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if setting != nil {
		data.ImageBuildSettings = setting.MustAsImageBuildSettings()
	}
	return nil
}

func (s *service) resetRepoCheckoutDir(
	data *repoDeploymentData,
) error {
	if err := os.RemoveAll(data.CheckoutDir); err != nil {
		return apperrors.Wrap(err)
	}
	if err := os.MkdirAll(data.CheckoutDir, base.DirModeDefault); err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (s *service) addCmdOutToLogs(
	ctx context.Context,
	msg string,
	isErr bool,
	logStore *tasklog.Store,
) {
	if logStore == nil || len(msg) == 0 {
		return
	}
	fn := gofn.If(isErr, tasklog.NewErrFrame, tasklog.NewOutFrame)
	_ = logStore.Add(ctx, fn(msg, tasklog.TsNow))
}
