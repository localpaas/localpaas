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

func (uc *ProjectUC) UpdateProjectSettings(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectdto.UpdateProjectSettingsReq,
) (*projectdto.UpdateProjectSettingsResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		data := &updateProjectSettingsData{}
		err := uc.loadProjectSettingsDataForUpdate(ctx, db, req, data)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingProjectData{}
		err = uc.preparePersistingProjectSettings(req, data, persistingData)
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
	Project *entity.Project
}

func (uc *ProjectUC) loadProjectSettingsDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *projectdto.UpdateProjectSettingsReq,
	data *updateProjectSettingsData,
) error {
	var targetTypes []base.SettingType
	switch {
	case req.EnvVars != nil:
		targetTypes = append(targetTypes, base.SettingTypeEnvVar)
	case req.Settings != nil:
		targetTypes = append(targetTypes, base.SettingTypeProject)
	}

	project, err := uc.projectRepo.GetByID(ctx, db, req.ProjectID,
		bunex.SelectFor("UPDATE OF project"),
		bunex.SelectRelation("Settings",
			bunex.SelectWhere("setting.type IN (?)", bunex.In(targetTypes)),
		),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Project = project

	return nil
}

func (uc *ProjectUC) preparePersistingProjectSettings(
	req *projectdto.UpdateProjectSettingsReq,
	data *updateProjectSettingsData,
	persistingData *persistingProjectData,
) error {
	timeNow := timeutil.NowUTC()
	project := data.Project

	if req.EnvVars != nil {
		setting := project.GetSettingByType(base.SettingTypeEnvVar)
		if setting == nil {
			setting = &entity.Setting{
				ID:        gofn.Must(ulid.NewStringULID()),
				ObjectID:  project.ID,
				Type:      base.SettingTypeEnvVar,
				Status:    base.SettingStatusActive,
				CreatedAt: timeNow,
			}
		}
		setting.UpdatedAt = timeNow
		err := setting.SetData(&entity.EnvVars{Data: gofn.MapSlice(req.EnvVars, func(v *projectdto.EnvVarReq) *entity.EnvVar {
			return v.ToEntity()
		})})
		if err != nil {
			return apperrors.Wrap(err)
		}
		persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
	}

	if req.Settings != nil {
		setting := project.GetSettingByType(base.SettingTypeProject)
		if setting == nil {
			setting = &entity.Setting{
				ID:        gofn.Must(ulid.NewStringULID()),
				ObjectID:  project.ID,
				Type:      base.SettingTypeProject,
				Status:    base.SettingStatusActive,
				CreatedAt: timeNow,
			}
		}
		setting.UpdatedAt = timeNow
		err := setting.SetData(req.Settings.ToEntity())
		if err != nil {
			return apperrors.Wrap(err)
		}
		persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
	}

	return nil
}
