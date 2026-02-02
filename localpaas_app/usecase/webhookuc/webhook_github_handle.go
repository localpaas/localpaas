package webhookuc

import (
	"context"
	"errors"

	"github.com/gitsight/go-vcsurl"
	"github.com/go-playground/webhooks/v6/github"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/service/appservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/webhookuc/webhookdto"
)

func (uc *WebhookUC) HandleWebhookGithub(
	ctx context.Context,
	req *webhookdto.HandleWebhookGithubReq,
) (*webhookdto.HandleWebhookGithubResp, error) {
	var persistingData *appservice.PersistingAppData
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		persistingData = &appservice.PersistingAppData{}
		mapWebhookSecret, err := uc.loadWebhookSecrets(ctx, db)
		if err != nil {
			return apperrors.Wrap(err)
		}

		var appsToRedeploy []*entity.App
		for secret, apps := range mapWebhookSecret {
			success, err := uc.processGithubWebhook(req, secret, apps, &appsToRedeploy)
			if err != nil {
				return apperrors.Wrap(err)
			}
			if success {
				break
			}
		}

		for _, app := range appsToRedeploy {
			err = uc.createAppDeployment(app, persistingData)
			if err != nil {
				return apperrors.Wrap(err)
			}
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

	return &webhookdto.HandleWebhookGithubResp{}, nil
}

func (uc *WebhookUC) processGithubWebhook(
	req *webhookdto.HandleWebhookGithubReq,
	secret string,
	apps []*entity.App,
	appsToRedeploy *[]*entity.App,
) (bool, error) {
	hook, err := github.New(github.Options.Secret(secret))
	if err != nil {
		return false, nil //nolint
	}
	payload, err := hook.Parse(req.Request, github.PushEvent)
	if err != nil {
		if errors.Is(err, github.ErrEventNotFound) { // ok event wasn't one of the ones asked to be parsed
			return true, nil
		}
		return false, nil //nolint
	}

	switch payload.(type) { //nolint
	case github.PushPayload:
		push := payload.(github.PushPayload) //nolint
		repoRef := push.Ref
		repoURL, err := vcsurl.Parse(push.Repository.HTMLURL)
		if err != nil {
			return false, apperrors.Wrap(err)
		}
		for _, app := range apps {
			if flag, _ := uc.shouldRedeployApp(app, repoURL, repoRef); flag {
				*appsToRedeploy = append(*appsToRedeploy, app)
			}
		}
	}
	return true, nil
}

func (uc *WebhookUC) loadWebhookSecrets(
	ctx context.Context,
	db database.IDB,
) (map[string][]*entity.App, error) {
	apps, _, err := uc.appRepo.List(ctx, db, "", nil,
		bunex.SelectFor("UPDATE OF app"),
		bunex.SelectColumns("id", "project_id", "parent_id", "webhook_secret"),
		bunex.SelectWhere("app.status = ?", base.AppStatusActive),
		bunex.SelectWhere("app.webhook_secret IS NOT NULL"),
		bunex.SelectJoin("JOIN projects ON projects.id = app.project_id"),
		bunex.SelectWhere("app.status = ?", base.AppStatusActive),
		bunex.SelectWhere("projects.status = ?", base.ProjectStatusActive),
		bunex.SelectRelation("Settings",
			bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
			bunex.SelectWhere("setting.type = ?", base.SettingTypeAppDeployment),
		),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	res := make(map[string][]*entity.App, len(apps))
	for _, app := range apps {
		if app.WebhookSecret == "" || len(app.Settings) == 0 {
			continue
		}
		res[app.WebhookSecret] = append(res[app.WebhookSecret], app)
	}
	return res, nil
}

func (uc *WebhookUC) shouldRedeployApp(
	app *entity.App,
	repoURL *vcsurl.VCS,
	repoRef string,
) (bool, error) {
	deploymentSettings, err := app.Settings[0].AsAppDeploymentSettings()
	if err != nil {
		return false, apperrors.Wrap(err)
	}

	repoSource := deploymentSettings.RepoSource
	if deploymentSettings.ActiveMethod != base.DeploymentMethodRepo || repoSource == nil {
		return false, nil
	}
	if repoRef != repoSource.RepoRef {
		return false, nil
	}
	url, err := vcsurl.Parse(repoSource.RepoURL)
	if err != nil {
		return false, apperrors.Wrap(err)
	}
	if url.ID != repoURL.ID {
		return false, nil
	}
	return true, nil
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
