package apppreviewserviceimpl

import (
	"context"
	"fmt"

	"github.com/moby/moby/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/githelper"
	"github.com/localpaas/localpaas/localpaas_app/service/appcopyservice"
	"github.com/localpaas/localpaas/localpaas_app/service/apppreviewservice"
)

type createPreviewData struct {
	*apppreviewservice.CreatePreviewReq

	Project    *entity.Project
	App        *entity.App
	PullRef    string // normalized pull ref
	PullNumber uint64

	PreviewApp         *entity.App
	Deployment         *entity.Deployment
	DeploymentTask     *entity.Task
	DeploymentSettings *entity.AppDeploymentSettings
}

func (s *service) CreatePreview(
	ctx context.Context,
	db database.Tx,
	req *apppreviewservice.CreatePreviewReq,
) (_ *apppreviewservice.CreatePreviewResp, err error) {
	data := &createPreviewData{
		CreatePreviewReq: req,
	}

	err = s.loadAppDataForCreatingPreview(ctx, db, data)
	if err != nil {
		return nil, apperrors.New(err)
	}

	copyResp, err := s.appCopyService.CopyApp(ctx, db, &appcopyservice.AppCopyReq{
		SrcProject:    data.Project,
		SrcApp:        data.App,
		TargetProject: data.Project,
		OnCopyApp: func(targetApp, srcApp *entity.App) error {
			return s.onCopyApp(targetApp, srcApp, data)
		},
		OnCopySetting: func(targetApp *entity.App, setting *entity.Setting) (*entity.Setting, error) {
			return s.onCopyAppSetting(setting, data)
		},
		OnCopyService: func(targetSvc, srcSvc *swarm.Service) error {
			return s.onCopyAppService(targetSvc, srcSvc, data)
		},
	})
	if err != nil {
		return nil, apperrors.New(err)
	}
	data.PreviewApp = copyResp.TargetApp

	err = s.createDeploymentAndTask(ctx, data)
	if err != nil {
		return nil, apperrors.New(err)
	}

	err = s.persistAppPreviewData(ctx, db, data)
	if err != nil {
		return nil, apperrors.New(err)
	}

	return &apppreviewservice.CreatePreviewResp{
		PreviewApp:     data.PreviewApp,
		Deployment:     data.Deployment,
		DeploymentTask: data.DeploymentTask,
		OnCleanup:      copyResp.OnCleanup,
	}, nil
}

func (s *service) loadAppDataForCreatingPreview(
	ctx context.Context,
	db database.IDB,
	data *createPreviewData,
) (err error) {
	data.PullRef, data.PullNumber, err = githelper.NormalizePullRef(data.PullRequest)
	if err != nil {
		return apperrors.New(err)
	}

	app, err := s.appService.LoadApp(ctx, db, data.ProjectID, data.AppID, true, true,
		bunex.SelectFor("UPDATE OF app"),
		bunex.SelectExcludeColumns(entity.AppDefaultExcludeColumns...),
		bunex.SelectRelation("Project",
			bunex.SelectExcludeColumns(entity.ProjectDefaultExcludeColumns...),
		),
		bunex.SelectRelation("Settings",
			bunex.SelectWhere("setting.type = ?", base.SettingTypeAppDeployment),
		),
	)
	if err != nil {
		return apperrors.New(err)
	}

	deploymentSetting := app.GetSettingByType(base.SettingTypeAppDeployment)
	if deploymentSetting == nil {
		return apperrors.NewNotFound("Deployment settings")
	}
	deploymentSettings := deploymentSetting.MustAsAppDeploymentSettings()
	if deploymentSettings.ActiveMethod != base.DeploymentMethodRepo || deploymentSettings.RepoSource == nil {
		return apperrors.New(apperrors.ErrDeploymentMethodRepoRequired)
	}

	data.App = app
	data.Project = app.Project
	return nil
}

func (s *service) onCopyApp(
	targetApp, _ *entity.App,
	data *createPreviewData,
) error {
	targetApp.Name = fmt.Sprintf("pr-%v", data.PullNumber)
	targetApp.Env = data.App.Env
	targetApp.Status = base.AppStatusActive
	targetApp.ParentID = data.App.ID // Preview app must be a child app of the current
	return nil
}

func (s *service) onCopyAppSetting(
	setting *entity.Setting,
	data *createPreviewData,
) (*entity.Setting, error) {
	switch setting.Type { //nolint:exhaustive
	case base.SettingTypeApp:
		return nil, nil
	case base.SettingTypeAppDeployment:
		return s.onCopyDeploymentSetting(setting, data)
	case base.SettingTypeAppFeatures:
		return nil, nil
	case base.SettingTypeAppHttp:
		return s.onCopyHttpSetting(setting, data)
	case base.SettingTypeConfigFile:
		return nil, nil
	case base.SettingTypeEnvVar:
		return nil, nil
	case base.SettingTypeHealthcheck:
		return nil, nil
	case base.SettingTypeSchedJob:
		return nil, nil
	case base.SettingTypeSecret:
		return nil, nil
	default:
		return nil, nil
	}
}

func (s *service) onCopyDeploymentSetting(
	setting *entity.Setting,
	data *createPreviewData,
) (*entity.Setting, error) {
	deploymentSettings := setting.MustAsAppDeploymentSettings()
	deploymentSettings.RepoSource.RepoRef = data.PullRef
	deploymentSettings.RepoSource.CommitHash = "" // unset target commit
	data.DeploymentSettings = deploymentSettings

	setting.MustSetData(deploymentSettings)
	return setting, nil
}

func (s *service) onCopyHttpSetting(
	setting *entity.Setting,
	data *createPreviewData,
) (*entity.Setting, error) {
	httpSettings := setting.MustAsAppHttpSettings()

	currDomains := httpSettings.Domains
	httpSettings.Domains = nil
	for _, domain := range currDomains {
		domain.Domain = fmt.Sprintf("pr-%v.%v", data.PullNumber, domain.Domain)
		// TODO: handle SSL cert
		httpSettings.Domains = append(httpSettings.Domains, domain)
	}

	setting.MustSetData(httpSettings)
	return setting, nil
}

func (s *service) onCopyAppService(
	targetSvc, _ *swarm.Service,
	_ *createPreviewData,
) error {
	targetSvcSpec := &targetSvc.Spec
	if targetSvcSpec.Mode.Replicated != nil {
		targetSvcSpec.Mode.Replicated.Replicas = new(uint64(1))
	}
	return nil
}

func (s *service) createDeploymentAndTask(
	_ context.Context,
	data *createPreviewData,
) (err error) {
	previewApp := data.PreviewApp
	deployment, deploymentTask, err := s.appDeploymentService.CreateDeploymentAndTask(
		previewApp, data.DeploymentSettings)
	if err != nil {
		return apperrors.New(err)
	}

	if data.OnInitDeployment != nil {
		if err = data.OnInitDeployment(deployment); err != nil {
			return apperrors.New(err)
		}
	}
	if data.OnDeploymentTask != nil {
		if err = data.OnDeploymentTask(deploymentTask); err != nil {
			return apperrors.New(err)
		}
	}

	data.Deployment = deployment
	data.DeploymentTask = deploymentTask
	return nil
}

func (s *service) persistAppPreviewData(
	ctx context.Context,
	db database.IDB,
	data *createPreviewData,
) (err error) {
	err = s.deploymentRepo.Upsert(ctx, db, data.Deployment,
		entity.DeploymentUpsertingConflictCols, entity.DeploymentUpsertingUpdateCols)
	if err != nil {
		return apperrors.New(err)
	}

	err = s.taskRepo.Upsert(ctx, db, data.DeploymentTask,
		entity.TaskUpsertingConflictCols, entity.TaskUpsertingUpdateCols)
	if err != nil {
		return apperrors.New(err)
	}

	return nil
}
