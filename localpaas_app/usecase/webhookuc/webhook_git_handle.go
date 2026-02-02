package webhookuc

import (
	"context"

	"github.com/gitsight/go-vcsurl"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/service/appservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/webhookuc/webhookdto"
)

func (uc *WebhookUC) HandleGitWebhook(
	ctx context.Context,
	req *webhookdto.HandleGitWebhookReq,
) (*webhookdto.HandleGitWebhookResp, error) {
	var persistingData *appservice.PersistingAppData
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		persistingData = &appservice.PersistingAppData{}

		webhookSettings, err := uc.loadWebhookSettings(ctx, db, req)
		if err != nil {
			return apperrors.Wrap(err)
		}
		if len(webhookSettings) == 0 { // no matching webhook setting found
			return nil
		}

		err = uc.processGitWebhook(ctx, db, req, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		err = uc.appService.PersistAppData(ctx, db, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}
		return nil
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	// Schedule deployment tasks
	for _, task := range persistingData.UpsertingTasks {
		_ = uc.taskQueue.ScheduleTask(ctx, task)
	}

	return &webhookdto.HandleGitWebhookResp{}, nil
}

type eventData struct {
	Push *pushEventData
}

type pushEventData struct {
	RepoRef string
	RepoURL string
}

func (uc *WebhookUC) processGitWebhook(
	ctx context.Context,
	db database.IDB,
	req *webhookdto.HandleGitWebhookReq,
	persistingData *appservice.PersistingAppData,
) (err error) {
	data := &eventData{}
	switch req.GitSource {
	case base.GitSourceGithub:
		err = uc.processGithubWebhook(req, data)
	case base.GitSourceGitlab, base.GitSourceGitlabCustom:
		err = uc.processGitlabWebhook(req, data)
	case base.GitSourceGitea:
		err = uc.processGiteaWebhook(req, data)
	case base.GitSourceBitbucket:
		err = uc.processBitbucketWebhook(req, data)
	default:
		return apperrors.New(apperrors.ErrUnsupported).
			WithMsgLog("git source %s not supported", req.GitSource)
	}
	if err != nil {
		return apperrors.Wrap(err)
	}

	if data.Push != nil {
		apps, err := uc.findAppsToRedeployByPushEvent(ctx, db, data.Push)
		if err != nil {
			return apperrors.Wrap(err)
		}
		for _, app := range apps {
			err = uc.createAppDeployment(app, persistingData)
			if err != nil {
				return apperrors.Wrap(err)
			}
		}
	}
	return nil
}

func (uc *WebhookUC) loadWebhookSettings(
	ctx context.Context,
	db database.IDB,
	req *webhookdto.HandleGitWebhookReq,
) ([]*entity.Setting, error) {
	settings, _, err := uc.settingRepo.List(ctx, db, nil,
		bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
		bunex.SelectWhereIn("setting.type IN (?)", base.SettingTypeWebhook, base.SettingTypeGithubApp),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	res := make([]*entity.Setting, 0, len(settings))
	for _, setting := range settings {
		switch setting.Type { //nolint:exhaustive
		case base.SettingTypeGithubApp:
			if setting.MustAsGithubApp().WebhookSecret == req.Secret {
				res = append(res, setting)
			}
		case base.SettingTypeWebhook:
			if setting.MustAsWebhook().Secret == req.Secret {
				res = append(res, setting)
			}
		}
	}
	return res, nil
}

func (uc *WebhookUC) findAppsToRedeployByPushEvent(
	ctx context.Context,
	db database.IDB,
	pushEvent *pushEventData,
) ([]*entity.App, error) {
	apps, _, err := uc.appRepo.List(ctx, db, "", nil,
		bunex.SelectFor("UPDATE OF app"),
		bunex.SelectWhere("app.status = ?", base.AppStatusActive),
		bunex.SelectJoin("JOIN projects ON projects.id = app.project_id"),
		bunex.SelectWhere("projects.status = ?", base.ProjectStatusActive),
		bunex.SelectRelation("Settings",
			bunex.SelectWhere("setting.type = ?", base.SettingTypeAppDeployment),
			bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
		),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if len(apps) == 0 {
		return nil, nil
	}

	inRepoURL, err := vcsurl.Parse(pushEvent.RepoURL)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	appsToRedeploy := make([]*entity.App, 0, len(apps))
	for _, app := range apps {
		matching, err := uc.shouldRedeployAppByPushEvent(app, inRepoURL, pushEvent.RepoRef)
		if err == nil && matching {
			appsToRedeploy = append(appsToRedeploy, app)
		}
	}

	return appsToRedeploy, nil
}

func (uc *WebhookUC) shouldRedeployAppByPushEvent(
	app *entity.App,
	inRepoURL *vcsurl.VCS,
	inRepoRef string,
) (bool, error) {
	if len(app.Settings) == 0 {
		return false, nil
	}
	deploymentSettings, err := app.Settings[0].AsAppDeploymentSettings()
	if err != nil {
		return false, apperrors.Wrap(err)
	}

	repoSource := deploymentSettings.RepoSource
	if deploymentSettings.ActiveMethod != base.DeploymentMethodRepo || repoSource == nil {
		return false, nil
	}
	if inRepoRef != repoSource.RepoRef {
		return false, nil
	}
	if inRepoURL.Raw == repoSource.RepoURL {
		return true, nil
	}
	url, err := vcsurl.Parse(repoSource.RepoURL)
	if err != nil {
		return false, apperrors.Wrap(err)
	}
	return url.ID == inRepoURL.ID, nil
}

func (uc *WebhookUC) createAppDeployment(
	app *entity.App,
	persistingData *appservice.PersistingAppData,
) error {
	deploymentSettings, err := app.Settings[0].AsAppDeploymentSettings()
	if err != nil {
		return apperrors.Wrap(err)
	}

	deployment, task, err := uc.appService.CreateDeployment(app, deploymentSettings)
	if err != nil {
		return apperrors.Wrap(err)
	}
	persistingData.UpsertingDeployments = append(persistingData.UpsertingDeployments, deployment)
	persistingData.UpsertingTasks = append(persistingData.UpsertingTasks, task)

	return nil
}
