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
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc/projectdto"
	"github.com/localpaas/localpaas/pkg/timeutil"
	"github.com/localpaas/localpaas/pkg/ulid"
)

func (uc *ProjectUC) UpdateProjectEnvVars(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectdto.UpdateProjectEnvVarsReq,
) (*projectdto.UpdateProjectEnvVarsResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		settingsData := &updateProjectEnvVarsData{}
		err := uc.loadProjectEnvVarsDataForUpdate(ctx, db, req, settingsData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingProjectData{}
		err = uc.preparePersistingProjectEnvVars(auth, req, settingsData, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &projectdto.UpdateProjectEnvVarsResp{}, nil
}

type updateProjectEnvVarsData struct {
	Project          *entity.Project
	ExistingSettings *entity.Setting
}

func (uc *ProjectUC) loadProjectEnvVarsDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *projectdto.UpdateProjectEnvVarsReq,
	data *updateProjectEnvVarsData,
) error {
	project, err := uc.projectRepo.GetByID(ctx, db, req.ProjectID,
		bunex.SelectRelation("AllSettings",
			// Filter by `project target`
			bunex.SelectWhere("setting.target_type = ?", base.SettingTargetEnvVar),
		),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Project = project

	if len(project.AllSettings) > 0 {
		data.ExistingSettings = project.AllSettings[0]
	}

	return nil
}

func (uc *ProjectUC) preparePersistingProjectEnvVars(
	auth *basedto.Auth,
	req *projectdto.UpdateProjectEnvVarsReq,
	data *updateProjectEnvVarsData,
	persistingData *persistingProjectData,
) error {
	timeNow := timeutil.NowUTC()
	project := data.Project
	settings := data.ExistingSettings
	if settings == nil {
		settings = &entity.Setting{
			ID:         gofn.Must(ulid.NewStringULID()),
			TargetType: base.SettingTargetEnvVar,
			TargetID:   project.ID,
			CreatedAt:  timeNow,
			CreatedBy:  auth.User.ID,
		}
	}

	settings.UpdatedAt = timeNow
	settings.UpdatedBy = auth.User.ID

	err := settings.SetData(&entity.ProjectEnvVars{Data: req.EnvVars})
	if err != nil {
		return apperrors.Wrap(err)
	}

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, settings)
	return nil
}
