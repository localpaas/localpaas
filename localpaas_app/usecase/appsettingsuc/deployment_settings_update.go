package appsettingsuc

import (
	"context"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appsettingsuc/appsettingsdto"
)

func (uc *UC) UpdateAppDeploymentSettings(
	ctx context.Context,
	auth *basedto.Auth,
	req *appsettingsdto.UpdateAppDeploymentSettingsReq,
) (*appsettingsdto.UpdateAppDeploymentSettingsResp, error) {
	var persistingData *persistingAppData
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		data := &updateAppDeploymentSettingsData{}
		err := uc.loadAppDeploymentSettingsForUpdate(ctx, db, req, data)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData = &persistingAppData{}
		err = uc.prepareUpdatingAppDeploymentSettings(ctx, auth, data, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		err = uc.persistData(ctx, db, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}
		return nil
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	err = uc.postTransactionAppDeploymentSettings(ctx, persistingData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appsettingsdto.UpdateAppDeploymentSettingsResp{}, nil
}

type updateAppDeploymentSettingsData struct {
	App                   *entity.App
	DeploymentSettings    *entity.Setting
	NewDeploymentSettings *entity.AppDeploymentSettings
	RegistryAuth          *entity.Setting
	Errors                []string // stores errors
	Warnings              []string // stores warnings

	RepoLinksToDelete []*entity.ResLink
	RepoLinksToAdd    []string
}

func (uc *UC) loadAppDeploymentSettingsForUpdate(
	ctx context.Context,
	db database.Tx,
	req *appsettingsdto.UpdateAppDeploymentSettingsReq,
	data *updateAppDeploymentSettingsData,
) error {
	app, err := uc.appService.LoadApp(ctx, db, req.ProjectID, req.AppID, true, true,
		bunex.SelectExcludeColumns(entity.AppDefaultExcludeColumns...),
		bunex.SelectFor("UPDATE OF app"),
		bunex.SelectRelation("Project",
			bunex.SelectExcludeColumns(entity.ProjectDefaultExcludeColumns...),
		),
		bunex.SelectRelation("Settings",
			bunex.SelectWhere("setting.type = ?", base.SettingTypeAppDeployment),
		),
		bunex.SelectRelation("DstResLinks",
			// NOTE: for now, we only need repo links
			bunex.SelectWhereIn("res_link.dst_type IN (?)", base.ResourceTypeRepo),
		),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.App = app
	data.DeploymentSettings = app.GetSettingByType(base.SettingTypeAppDeployment)

	deploymentSettings := data.DeploymentSettings
	if deploymentSettings != nil && deploymentSettings.UpdateVer != req.UpdateVer {
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}

	newDeploymentSettings, err := req.ToEntity()
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.NewDeploymentSettings = newDeploymentSettings

	// Make sure all reference settings used in this settings exist actively
	_, err = uc.settingService.LoadReferenceObjectsByIDs(ctx, db, app.GetSettingScope(),
		true, true, newDeploymentSettings.GetRefObjectIDs())
	if err != nil {
		return apperrors.Wrap(err)
	}

	if newDeploymentSettings.RepoSource != nil {
		// When the cluster has multiple nodes, the result image must be pushed to a registry
		// that can be accessed by all the nodes in the cluster.
		isMultiNode, err := uc.clusterService.IsMultiNode(ctx)
		if err != nil {
			return apperrors.Wrap(err)
		}
		if isMultiNode && newDeploymentSettings.RepoSource.PushToRegistry.ID == "" {
			return apperrors.Wrap(apperrors.ErrMultiNodeClusterRequireRegistryForImages)
		}
	}

	// Calculate repo links update
	mapCurrentRepoLinks := make(map[string]*entity.ResLink, len(app.DstResLinks))
	for _, repoLink := range app.DstResLinks {
		if repoLink.DstType != base.ResourceTypeRepo {
			continue
		}
		mapCurrentRepoLinks[repoLink.DstID] = repoLink
	}
	if newDeploymentSettings.ActiveMethod == base.DeploymentMethodRepo {
		repoID := newDeploymentSettings.RepoSource.RepoID
		repoLink := mapCurrentRepoLinks[repoID]
		if repoLink == nil {
			data.RepoLinksToAdd = append(data.RepoLinksToAdd, repoID)
		} else {
			delete(mapCurrentRepoLinks, repoID)
		}
	}
	for _, repoLink := range mapCurrentRepoLinks {
		data.RepoLinksToDelete = append(data.RepoLinksToDelete, repoLink)
	}

	return nil
}

func (uc *UC) prepareUpdatingAppDeploymentSettings(
	_ context.Context,
	auth *basedto.Auth,
	data *updateAppDeploymentSettingsData,
	persistingData *persistingAppData,
) error {
	app := data.App
	setting := data.DeploymentSettings
	timeNow := timeutil.NowUTC()

	if setting == nil {
		setting = &entity.Setting{
			ID:        gofn.Must(ulid.NewStringULID()),
			Scope:     base.ObjectScopeApp,
			ObjectID:  app.ID,
			Type:      base.SettingTypeAppDeployment,
			CreatedAt: timeNow,
			Version:   entity.CurrentAppDeploymentSettingsVersion,
		}
		data.DeploymentSettings = setting
	}
	setting.UpdateVer++
	setting.UpdatedAt = timeNow
	setting.ExpireAt = time.Time{}
	setting.Status = base.SettingStatusActive
	setting.MustSetData(data.NewDeploymentSettings)
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)

	// Create a deployment and a task for it
	deployment, deploymentTask, err := uc.appDeploymentService.CreateDeploymentAndTask(app, data.NewDeploymentSettings)
	if err != nil {
		return apperrors.Wrap(err)
	}
	// Set trigger for the deployment
	deployment.Trigger = &entity.AppDeploymentTrigger{
		Source:   base.DeploymentTriggerSourceUser,
		SourceID: auth.User.ID,
	}

	persistingData.UpsertingDeployments = append(persistingData.UpsertingDeployments, deployment)
	persistingData.UpsertingTasks = append(persistingData.UpsertingTasks, deploymentTask)

	// Repo links
	for _, repoLink := range data.RepoLinksToDelete {
		repoLink.DeletedAt = timeNow
		persistingData.UpsertingResLinks = append(persistingData.UpsertingResLinks, repoLink)
	}
	for _, repoID := range data.RepoLinksToAdd {
		persistingData.UpsertingResLinks = append(persistingData.UpsertingResLinks, &entity.ResLink{
			SrcType: base.ResourceTypeApp,
			SrcID:   app.ID,
			DstType: base.ResourceTypeRepo,
			DstID:   repoID,
		})
	}

	return nil
}

func (uc *UC) postTransactionAppDeploymentSettings(
	ctx context.Context,
	persistingData *persistingAppData,
) error {
	for _, task := range persistingData.UpsertingTasks {
		err := uc.taskQueue.ScheduleTask(ctx, task)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}
	return nil
}
