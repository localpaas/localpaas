package projectenvuc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/service/projectservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectenvuc/projectenvdto"
	"github.com/localpaas/localpaas/pkg/timeutil"
	"github.com/localpaas/localpaas/pkg/ulid"
)

func (uc *ProjectEnvUC) UpdateProjectEnvSettings(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectenvdto.UpdateProjectEnvSettingsReq,
) (*projectenvdto.UpdateProjectEnvSettingsResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		settingsData := &updateProjectEnvSettingsData{}
		err := uc.loadProjectEnvSettingsDataForUpdate(ctx, db, req, settingsData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &projectservice.PersistingProjectData{}
		err = uc.preparePersistingProjectEnvSettings(auth, req, settingsData, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		return uc.projectService.PersistProjectData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &projectenvdto.UpdateProjectEnvSettingsResp{}, nil
}

type updateProjectEnvSettingsData struct {
	ProjectEnv       *entity.ProjectEnv
	ExistingSettings *entity.Setting
}

func (uc *ProjectEnvUC) loadProjectEnvSettingsDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *projectenvdto.UpdateProjectEnvSettingsReq,
	data *updateProjectEnvSettingsData,
) error {
	projectEnv, err := uc.projectEnvRepo.GetByID(ctx, db, req.ProjectEnvID,
		bunex.SelectRelation("MainSettings"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if projectEnv.ProjectID != req.ProjectID {
		return apperrors.New(apperrors.ErrUnauthorized)
	}
	data.ProjectEnv = projectEnv

	if len(projectEnv.MainSettings) > 0 {
		data.ExistingSettings = projectEnv.MainSettings[0]
	}

	return nil
}

func (uc *ProjectEnvUC) preparePersistingProjectEnvSettings(
	auth *basedto.Auth,
	req *projectenvdto.UpdateProjectEnvSettingsReq,
	data *updateProjectEnvSettingsData,
	persistingData *projectservice.PersistingProjectData,
) error {
	timeNow := timeutil.NowUTC()
	projectEnv := data.ProjectEnv
	settings := data.ExistingSettings
	if settings == nil {
		settings = &entity.Setting{
			ID:         gofn.Must(ulid.NewStringULID()),
			TargetType: base.SettingTargetProjectEnv,
			TargetID:   projectEnv.ID,
			CreatedAt:  timeNow,
			CreatedBy:  auth.User.ID,
		}
	}

	settings.UpdatedAt = timeNow
	settings.UpdatedBy = auth.User.ID

	var settingsData *entity.ProjectEnvSettings

	// Do a copy fields to fields
	err := copier.Copy(&settingsData, req.Settings)
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = settings.SetData(settingsData)
	if err != nil {
		return apperrors.Wrap(err)
	}

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, settings)
	return nil
}
