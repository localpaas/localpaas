package appuc

import (
	"context"
	"errors"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/slugify"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/service/appservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
	"github.com/localpaas/localpaas/services/docker"
)

const (
	appKeyMaxLen = 100
)

func (uc *AppUC) CreateApp(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.CreateAppReq,
) (*appdto.CreateAppResp, error) {
	var appData *createAppData
	var persistingData *persistingAppData
	var createdApp *entity.App
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		appData = &createAppData{}
		err := uc.loadAppData(ctx, db, req, appData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData = &persistingAppData{}
		err = uc.preparePersistingApp(ctx, req, appData, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		createdApp = persistingData.UpsertingApps[0]

		// Create a service in docker for the app
		res, err := uc.dockerManager.ServiceCreate(ctx, gofn.Must(appData.ServiceSpec.ToSwarmServiceSpec()))
		if err != nil {
			return apperrors.NewInfra(err)
		}
		if res.ID == "" { // should never happen
			return apperrors.New(apperrors.ErrInfraInternal).
				WithNTParam("Error", "empty service ID returned")
		}
		createdApp.ServiceID = res.ID

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.CreateAppResp{
		Data: &basedto.ObjectIDResp{ID: createdApp.ID},
	}, nil
}

type createAppData struct {
	Project     *entity.Project
	AppKey      string
	ServiceSpec *docker.ServiceSpec
}

func (uc *AppUC) loadAppData(
	ctx context.Context,
	db database.IDB,
	req *appdto.CreateAppReq,
	data *createAppData,
) error {
	project, err := uc.projectRepo.GetByID(ctx, db, req.ProjectID,
		bunex.SelectFor("UPDATE OF project"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if project.Status != base.ProjectStatusActive {
		return apperrors.Wrap(apperrors.ErrResourceInactive)
	}
	data.Project = project

	data.AppKey = project.Key + "_" + slugify.SlugifyEx(req.Name, []string{"-", "_"}, appKeyMaxLen)

	app, err := uc.appRepo.GetByKey(ctx, db, project.ID, data.AppKey)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if app != nil {
		return apperrors.NewAlreadyExist("App").
			WithMsgLog("app '%s' already exists", data.AppKey)
	}

	return nil
}

type persistingAppData struct {
	appservice.PersistingAppData
}

func (uc *AppUC) preparePersistingApp(
	ctx context.Context,
	req *appdto.CreateAppReq,
	data *createAppData,
	persistingData *persistingAppData,
) error {
	timeNow := timeutil.NowUTC()
	project := data.Project
	app := &entity.App{
		ID:        gofn.Must(ulid.NewStringULID()),
		ProjectID: project.ID,
		Key:       data.AppKey,
		CreatedAt: timeNow,
	}

	uc.preparePersistingAppBase(app, req.AppBaseReq, timeNow, persistingData)
	uc.preparePersistingAppTags(app, req.Tags, 0, persistingData)
	err := uc.preparePersistingAppSettingsDefault(ctx, app, timeNow, data, persistingData)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (uc *AppUC) preparePersistingAppBase(
	app *entity.App,
	req *appdto.AppBaseReq,
	timeNow time.Time,
	persistingData *persistingAppData,
) {
	app.Name = req.Name
	app.Status = req.Status
	app.Note = req.Note
	app.UpdatedAt = timeNow

	persistingData.UpsertingApps = append(persistingData.UpsertingApps, app)
}

func (uc *AppUC) preparePersistingAppTags(
	app *entity.App,
	tags []string,
	startDisplayOrder int,
	persistingData *persistingAppData,
) {
	displayOrder := startDisplayOrder
	for _, tag := range tags {
		persistingData.UpsertingTags = append(persistingData.UpsertingTags,
			&entity.AppTag{
				AppID:        app.ID,
				Tag:          tag,
				DisplayOrder: displayOrder,
			})
		displayOrder++
	}
}

func (uc *AppUC) preparePersistingAppSettingsDefault(
	ctx context.Context,
	app *entity.App,
	timeNow time.Time,
	data *createAppData,
	persistingData *persistingAppData,
) error {
	setting := &entity.Setting{
		ID:        gofn.Must(ulid.NewStringULID()),
		ObjectID:  app.ID,
		Type:      base.SettingTypeServiceSpec,
		Status:    base.SettingStatusActive,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	serviceSpec := &docker.ServiceSpec{
		Name: app.Key,
		Labels: map[string]string{
			docker.StackLabelNamespace: data.Project.Key,
		},
		Image:       "crccheck/hello-world:latest", // TODO: test image
		ServiceMode: docker.ServiceModeReplicated,
		Replicas:    1,
		Hostname:    app.Key,
		Networks: []*docker.NetworkAttachment{
			{
				Target: data.Project.GetDefaultNetworkName(),
			},
		},
	}
	setting.MustSetData(serviceSpec)
	data.ServiceSpec = serviceSpec

	httpSettings := &entity.AppHttpSettings{}
	err := uc.nginxService.InitAppConfig(ctx, app, httpSettings)
	if err != nil {
		return apperrors.Wrap(err)
	}

	httpSettingsObj := &entity.Setting{
		ID:        gofn.Must(ulid.NewStringULID()),
		Type:      base.SettingTypeAppHttp,
		Status:    base.SettingStatusActive,
		ObjectID:  app.ID,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}
	httpSettingsObj.MustSetData(httpSettings)
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, httpSettingsObj)

	return nil
}

func (uc *AppUC) persistData(
	ctx context.Context,
	db database.IDB,
	persistingData *persistingAppData,
) error {
	err := uc.appService.PersistAppData(ctx, db, &persistingData.PersistingAppData)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
