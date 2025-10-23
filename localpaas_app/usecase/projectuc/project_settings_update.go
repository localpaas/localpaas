package projectuc

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
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc/projectdto"
	"github.com/localpaas/localpaas/pkg/timeutil"
	"github.com/localpaas/localpaas/pkg/ulid"
)

func (uc *ProjectUC) UpdateProjectSettings(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectdto.UpdateProjectSettingsReq,
) (*projectdto.UpdateProjectSettingsResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		settingsData := &updateProjectSettingsData{}
		err := uc.loadProjectSettingsDataForUpdate(ctx, db, req, settingsData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingProjectData{}
		err = uc.preparePersistingProjectSettings(req, settingsData, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &projectdto.UpdateProjectSettingsResp{}, nil
}

type updateProjectSettingsData struct {
	Project          *entity.Project
	ExistingSettings *entity.Setting
}

func (uc *ProjectUC) loadProjectSettingsDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *projectdto.UpdateProjectSettingsReq,
	data *updateProjectSettingsData,
) error {
	project, err := uc.projectRepo.GetByID(ctx, db, req.ProjectID,
		bunex.SelectFor("UPDATE OF project"),
		bunex.SelectRelation("MainSettings"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Project = project

	if len(project.MainSettings) > 0 {
		data.ExistingSettings = project.MainSettings[0]
	}

	return nil
}

func (uc *ProjectUC) preparePersistingProjectSettings(
	req *projectdto.UpdateProjectSettingsReq,
	data *updateProjectSettingsData,
	persistingData *persistingProjectData,
) error {
	timeNow := timeutil.NowUTC()
	project := data.Project
	settings := data.ExistingSettings
	if settings == nil {
		settings = &entity.Setting{
			ID:        gofn.Must(ulid.NewStringULID()),
			Type:      base.SettingTypeProject,
			ObjectID:  project.ID,
			CreatedAt: timeNow,
		}
	}

	settings.UpdatedAt = timeNow

	var settingsData *entity.ProjectSettings

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
