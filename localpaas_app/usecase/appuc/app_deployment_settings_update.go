package appuc

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
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/pkg/githelper"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

func (uc *AppUC) UpdateAppDeploymentSettings(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.UpdateAppDeploymentSettingsReq,
) (*appdto.UpdateAppDeploymentSettingsResp, error) {
	var data *updateAppDeploymentSettingsData
	var persistingData *persistingAppData
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		data = &updateAppDeploymentSettingsData{}
		err := uc.loadAppDeploymentSettingsForUpdate(ctx, db, req, data)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData = &persistingAppData{}
		err = uc.prepareUpdatingAppDeploymentSettings(ctx, db, req, data, persistingData)
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

	return &appdto.UpdateAppDeploymentSettingsResp{}, nil
}

type updateAppDeploymentSettingsData struct {
	App                    *entity.App
	DeploymentSettings     *entity.Setting
	CurrDeploymentSettings *entity.AppDeploymentSettings
	RegistryAuth           *entity.Setting
	Errors                 []string // stores errors
	Warnings               []string // stores warnings
}

func (uc *AppUC) loadAppDeploymentSettingsForUpdate(
	ctx context.Context,
	db database.Tx,
	req *appdto.UpdateAppDeploymentSettingsReq,
	data *updateAppDeploymentSettingsData,
) error {
	app, err := uc.appService.LoadApp(ctx, db, req.ProjectID, req.AppID, true, true,
		bunex.SelectExcludeColumns(entity.AppDefaultExcludeColumns...),
		bunex.SelectFor("UPDATE OF app"),
		bunex.SelectRelation("Project"),
		bunex.SelectRelation("Settings",
			bunex.SelectWhere("setting.type = ?", base.SettingTypeAppDeployment),
		),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.App = app
	data.DeploymentSettings, _ = gofn.First(app.Settings)

	deploymentSettings := data.DeploymentSettings
	if deploymentSettings != nil && deploymentSettings.UpdateVer != req.UpdateVer {
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}

	// Parse the current deployment settings
	if deploymentSettings != nil && deploymentSettings.IsActive() {
		settingData, err := deploymentSettings.AsAppDeploymentSettings()
		if err != nil {
			return apperrors.New(err).WithMsgLog("failed to parse app deployment settings")
		}
		data.CurrDeploymentSettings = settingData
	}

	if req.RepoSource != nil {
		// Normalize repo type (currently supports git type only)
		if req.RepoSource.RepoType == base.RepoTypeGit {
			req.RepoSource.RepoRef = string(githelper.NormalizeRepoRef(req.RepoSource.RepoRef))
		}

		// When the cluster has multiple nodes, the result image must be pushed to a registry
		// that can be accessed by all the nodes in the cluster.
		isMultiNode, err := uc.clusterService.IsMultiNode(ctx)
		if err != nil {
			return apperrors.Wrap(err)
		}
		if isMultiNode && req.RepoSource.PushToRegistry.ID == "" {
			return apperrors.Wrap(apperrors.ErrMultiNodeClusterRequireRegistryForImages)
		}
	}

	return nil
}

func (uc *AppUC) prepareUpdatingAppDeploymentSettings(
	ctx context.Context,
	db database.IDB,
	req *appdto.UpdateAppDeploymentSettingsReq,
	data *updateAppDeploymentSettingsData,
	persistingData *persistingAppData,
) error {
	app := data.App
	setting := data.DeploymentSettings
	timeNow := timeutil.NowUTC()

	if setting == nil {
		setting = &entity.Setting{
			ID:        gofn.Must(ulid.NewStringULID()),
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

	newDeploymentSettings, err := uc.buildNewAppDeploymentSettings(req, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Validation: Make sure all reference settings used in this deployment settings exist actively
	_, err = uc.settingService.LoadReferenceObjects(ctx, db, base.SettingScopeApp, app.ID, app.ProjectID,
		true, true, setting)
	if err != nil {
		return apperrors.Wrap(err)
	}

	setting.MustSetData(newDeploymentSettings)
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)

	// Create a deployment and a task for it
	deployment, deploymentTask, err := uc.appService.CreateDeployment(app, newDeploymentSettings)
	if err != nil {
		return apperrors.Wrap(err)
	}
	// Set trigger for the deployment
	deployment.Trigger = &entity.AppDeploymentTrigger{
		Source: base.DeploymentTriggerSourceUser,
	}

	persistingData.UpsertingDeployments = append(persistingData.UpsertingDeployments, deployment)
	persistingData.UpsertingTasks = append(persistingData.UpsertingTasks, deploymentTask)

	return nil
}

func (uc *AppUC) buildNewAppDeploymentSettings(
	req *appdto.UpdateAppDeploymentSettingsReq,
	data *updateAppDeploymentSettingsData,
) (*entity.AppDeploymentSettings, error) {
	newDeploymentSettings := data.CurrDeploymentSettings
	if newDeploymentSettings == nil {
		newDeploymentSettings = &entity.AppDeploymentSettings{}
	}

	newDeploymentSettings.ActiveMethod = req.ActiveMethod
	if req.ImageSource != nil {
		err := copier.Copy(&newDeploymentSettings.ImageSource, req.ImageSource)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
	}
	if req.RepoSource != nil {
		err := copier.Copy(&newDeploymentSettings.RepoSource, req.RepoSource)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
	}

	if req.PreDeploymentCommand != nil {
		newDeploymentSettings.PreDeploymentCommand = req.PreDeploymentCommand
	}
	if req.PostDeploymentCommand != nil {
		newDeploymentSettings.PostDeploymentCommand = req.PostDeploymentCommand
	}

	return newDeploymentSettings, nil
}

func (uc *AppUC) postTransactionAppDeploymentSettings(
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
