package projectserviceimpl

import (
	"context"
	"errors"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
)

func (s *service) InitRootProject(
	ctx context.Context,
	db database.IDB,
) error {
	project, err := s.projectRepo.GetByKey(ctx, db, base.LocalpaasProjectKey)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if project == nil {
		timeNow := timeutil.NowUTC()
		project = &entity.Project{
			ID:        gofn.Must(ulid.NewStringULID()),
			Name:      base.LocalpaasProjectName,
			Key:       base.LocalpaasProjectKey,
			Status:    base.ProjectStatusActive,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		}

		// Get admin account and assign it to project as owner
		users, _, err := s.userRepo.List(ctx, db, nil,
			bunex.SelectColumns("id"),
			bunex.SelectWhere("role = ?", base.UserRoleAdmin),
			bunex.SelectWhere("status = ?", base.UserStatusActive),
			bunex.SelectOrder("created_at"),
			bunex.SelectLimit(1),
		)
		if err != nil {
			return apperrors.Wrap(err)
		}
		if len(users) > 0 {
			project.OwnerID = users[0].ID
		}
	}

	err = s.projectRepo.Upsert(ctx, db, project,
		entity.ProjectUpsertingConflictCols, entity.ProjectUpsertingUpdateCols)
	if err != nil {
		return apperrors.Wrap(err)
	}

	newApps, _, err := s.SyncProject(ctx, db, project)
	if err != nil {
		return apperrors.Wrap(err)
	}

	for _, app := range newApps {
		if app.Key == base.LocalpaasAppKey {
			err = s.initRootAppLocalpaas(ctx, db, app)
			if err != nil {
				return apperrors.Wrap(err)
			}
		}
	}

	return nil
}

func (s *service) initRootAppLocalpaas(
	ctx context.Context,
	db database.IDB,
	app *entity.App,
) error {
	timeNow := timeutil.NowUTC()
	cfg := config.Current

	// Add HTTP settings for the main app
	dbHttpSetting := &entity.Setting{
		ID:        gofn.Must(ulid.NewStringULID()),
		Scope:     base.SettingScopeApp,
		ObjectID:  app.ID,
		Type:      base.SettingTypeAppHttp,
		Status:    base.SettingStatusActive,
		Version:   entity.CurrentAppHttpSettingsVersion,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}
	httpSettings := &entity.AppHttpSettings{
		ExposePublicly: true,
		Domains: []*entity.AppDomain{
			{
				Enabled:       true,
				Domain:        cfg.AppDomain,
				ContainerPort: cfg.HTTPServer.Port,
				ForceHttps:    true,
			},
		},
	}
	dbHttpSetting.MustSetData(httpSettings)

	err := s.settingRepo.Upsert(ctx, db, dbHttpSetting,
		entity.SettingUpsertingConflictCols, entity.SettingUpsertingUpdateCols)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
