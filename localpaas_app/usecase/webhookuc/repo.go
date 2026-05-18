package webhookuc

import (
	"context"
	"time"

	"github.com/gitsight/go-vcsurl"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/service/appservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/webhookuc/webhookdto"
)

func (uc *UC) HandleRepoWebhook(
	ctx context.Context,
	req *webhookdto.HandleRepoWebhookReq,
) (*webhookdto.HandleRepoWebhookResp, error) {
	var persistingData *appservice.PersistingAppData
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		data := &handleRepoWebhookData{}
		persistingData = &appservice.PersistingAppData{}

		err := uc.loadWebhookSettings(ctx, db, req, data)
		if err != nil {
			return apperrors.Wrap(err)
		}

		err = uc.processRepoWebhook(ctx, db, req, data, persistingData)
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

	return &webhookdto.HandleRepoWebhookResp{}, nil
}

type handleRepoWebhookData struct {
	RepoWebhook *entity.RepoWebhook
}

type repoEventData struct {
	Push *repoPushEventData
}

type repoPushEventData struct {
	RepoRef  string
	RepoURL  string
	ChangeID string

	parsedURL *vcsurl.VCS
}

func (uc *UC) processRepoWebhook(
	ctx context.Context,
	db database.IDB,
	req *webhookdto.HandleRepoWebhookReq,
	data *handleRepoWebhookData,
	persistingData *appservice.PersistingAppData,
) (err error) {
	eventData := &repoEventData{}
	switch data.RepoWebhook.Kind {
	case base.WebhookKindGithub:
		err = uc.processGithubWebhook(req, eventData)
	case base.WebhookKindGitlab:
		err = uc.processGitlabWebhook(req, eventData)
	case base.WebhookKindGitea:
		err = uc.processGiteaWebhook(req, eventData)
	case base.WebhookKindBitbucket:
		err = uc.processBitbucketWebhook(req, eventData)
	case base.WebhookKindGogs:
		err = uc.processGogsWebhook(req, eventData)
	case base.WebhookKindAzureDevOps:
		err = uc.processAzureDevOpsWebhook(req, eventData)
	default:
		return apperrors.New(apperrors.ErrUnsupported).
			WithMsgLog("webhook kind '%s' not supported", data.RepoWebhook.Kind)
	}
	if err != nil {
		return apperrors.Wrap(err)
	}

	if eventData.Push != nil {
		apps, err := uc.findAppsToRedeployByPushEvent(ctx, db, eventData.Push)
		if err != nil {
			return apperrors.Wrap(err)
		}
		for _, app := range apps {
			err = uc.createAppDeploymentByPushEvent(app, eventData.Push, persistingData)
			if err != nil {
				return apperrors.Wrap(err)
			}
		}
	}
	return nil
}

func (uc *UC) loadWebhookSettings(
	ctx context.Context,
	db database.IDB,
	req *webhookdto.HandleRepoWebhookReq,
	data *handleRepoWebhookData,
) error {
	settings, _, err := uc.settingRepo.List(ctx, db, nil, nil,
		bunex.SelectWhere("setting.id = ?", req.ID),
		bunex.SelectWhere("setting.status = ?", base.SettingStatusActive),
		bunex.SelectWhereIn("setting.type IN (?)", base.SettingTypeRepoWebhook, base.SettingTypeGithubApp),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}

	for _, setting := range settings {
		repoWebhook := setting.MustAsRepoWebhook()
		if repoWebhook.Secret != req.Secret {
			continue
		}
		data.RepoWebhook = repoWebhook
	}
	if data.RepoWebhook == nil {
		return apperrors.NewNotFound("Repo webhook settings")
	}

	return nil
}

func (uc *UC) findAppsToRedeployByPushEvent(
	ctx context.Context,
	db database.IDB,
	pushEvent *repoPushEventData,
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

	pushEvent.parsedURL, err = vcsurl.Parse(pushEvent.RepoURL)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	appsToRedeploy := make([]*entity.App, 0, len(apps))
	for _, app := range apps {
		matching, err := uc.shouldRedeployAppByPushEvent(ctx, db, app, pushEvent)
		if err == nil && matching {
			appsToRedeploy = append(appsToRedeploy, app)
		}
	}

	return appsToRedeploy, nil
}

func (uc *UC) shouldRedeployAppByPushEvent(
	ctx context.Context,
	db database.IDB,
	app *entity.App,
	pushEvent *repoPushEventData,
) (bool, error) {
	if len(app.Settings) == 0 {
		return false, nil
	}
	deploymentSettings, err := app.Settings[0].AsAppDeploymentSettings()
	if err != nil {
		return false, apperrors.Wrap(err)
	}

	shouldRedeploy := false
	for {
		repoSource := deploymentSettings.RepoSource
		if deploymentSettings.ActiveMethod != base.DeploymentMethodRepo || repoSource == nil {
			break
		}
		if pushEvent.RepoRef != repoSource.RepoRef {
			break
		}
		if pushEvent.RepoURL == repoSource.RepoURL {
			shouldRedeploy = true
			break
		}
		parsedURL, err := vcsurl.Parse(repoSource.RepoURL)
		if err != nil {
			return false, apperrors.Wrap(err)
		}
		shouldRedeploy = parsedURL.ID == pushEvent.parsedURL.ID
		break //nolint:staticcheck
	}

	if !shouldRedeploy {
		return false, nil
	}

	// Make sure there is no duplicated deployment having the same `change id`
	if pushEvent.ChangeID == "" {
		return true, nil
	}
	deployments, _, err := uc.deploymentRepo.List(ctx, db, app.ID, nil,
		bunex.SelectColumns("id"),
		bunex.SelectLimit(1),
		bunex.SelectWhere("deployment.created_at > ?", timeutil.NowUTC().Add(-time.Minute)),
		bunex.SelectWhere("deployment.trigger->>'source' = ?", base.DeploymentTriggerSourceRepoWebhook),
		bunex.SelectWhere("deployment.trigger->>'id' = ?", pushEvent.ChangeID),
	)
	if err != nil {
		return false, apperrors.Wrap(err)
	}
	return len(deployments) == 0, nil
}

func (uc *UC) createAppDeploymentByPushEvent(
	app *entity.App,
	pushEvent *repoPushEventData,
	persistingData *appservice.PersistingAppData,
) error {
	deploymentSettings, err := app.Settings[0].AsAppDeploymentSettings()
	if err != nil {
		return apperrors.Wrap(err)
	}

	deployment, task, err := uc.appDeploymentService.CreateDeploymentAndTask(app, deploymentSettings)
	if err != nil {
		return apperrors.Wrap(err)
	}
	// Set trigger for the deployment
	deployment.Trigger = &entity.AppDeploymentTrigger{
		Source: base.DeploymentTriggerSourceRepoWebhook,
		ID:     pushEvent.ChangeID,
	}

	persistingData.UpsertingDeployments = append(persistingData.UpsertingDeployments, deployment)
	persistingData.UpsertingTasks = append(persistingData.UpsertingTasks, task)
	return nil
}
