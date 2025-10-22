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
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/service/projectservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectenvuc/projectenvdto"
	"github.com/localpaas/localpaas/pkg/timeutil"
	"github.com/localpaas/localpaas/pkg/ulid"
)

func (uc *ProjectEnvUC) UpdateProjectEnvEnvVars(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectenvdto.UpdateProjectEnvEnvVarsReq,
) (*projectenvdto.UpdateProjectEnvEnvVarsResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		envData := &updateProjectEnvEnvVarsData{}
		err := uc.loadProjectEnvEnvVarsDataForUpdate(ctx, db, req, envData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &projectservice.PersistingProjectData{}
		err = uc.preparePersistingProjectEnvEnvVars(auth, req, envData, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		return uc.projectService.PersistProjectData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &projectenvdto.UpdateProjectEnvEnvVarsResp{}, nil
}

type updateProjectEnvEnvVarsData struct {
	ProjectEnv       *entity.ProjectEnv
	ExistingSettings *entity.Setting
}

func (uc *ProjectEnvUC) loadProjectEnvEnvVarsDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *projectenvdto.UpdateProjectEnvEnvVarsReq,
	data *updateProjectEnvEnvVarsData,
) error {
	projectEnv, err := uc.projectEnvRepo.GetByID(ctx, db, req.ProjectEnvID,
		bunex.SelectRelation("EnvVarsSettings"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if projectEnv.ProjectID != req.ProjectID {
		return apperrors.New(apperrors.ErrUnauthorized)
	}
	data.ProjectEnv = projectEnv

	if len(projectEnv.EnvVarsSettings) > 0 {
		data.ExistingSettings = projectEnv.EnvVarsSettings[0]
	}

	return nil
}

func (uc *ProjectEnvUC) preparePersistingProjectEnvEnvVars(
	auth *basedto.Auth,
	req *projectenvdto.UpdateProjectEnvEnvVarsReq,
	data *updateProjectEnvEnvVarsData,
	persistingData *projectservice.PersistingProjectData,
) error {
	timeNow := timeutil.NowUTC()
	projectEnv := data.ProjectEnv
	settings := data.ExistingSettings
	if settings == nil {
		settings = &entity.Setting{
			ID:         gofn.Must(ulid.NewStringULID()),
			TargetType: base.SettingTargetEnvVar,
			TargetID:   projectEnv.ID,
			CreatedAt:  timeNow,
			CreatedBy:  auth.User.ID,
		}
	}

	settings.UpdatedAt = timeNow
	settings.UpdatedBy = auth.User.ID

	err := settings.SetData(&entity.ProjectEnvEnvVars{Data: req.EnvVars})
	if err != nil {
		return apperrors.Wrap(err)
	}

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, settings)
	return nil
}
