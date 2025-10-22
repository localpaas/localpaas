package projectenvuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/service/projectservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectenvuc/projectenvdto"
	"github.com/localpaas/localpaas/pkg/timeutil"
)

func (uc *ProjectEnvUC) DeleteProjectEnv(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectenvdto.DeleteProjectEnvReq,
) (*projectenvdto.DeleteProjectEnvResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		envData := &deleteProjectEnvData{}
		err := uc.loadProjectEnvDataForDelete(ctx, db, req, envData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &projectservice.PersistingProjectData{}
		uc.prepareDeletingProjectEnv(envData, persistingData)

		// TODO: delete all apps in the deleted env

		return uc.projectService.PersistProjectData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &projectenvdto.DeleteProjectEnvResp{}, nil
}

type deleteProjectEnvData struct {
	Project            *entity.Project
	DeletingProjectEnv *entity.ProjectEnv
	UpdatingOrderEnvs  []*entity.ProjectEnv
}

func (uc *ProjectEnvUC) loadProjectEnvDataForDelete(
	ctx context.Context,
	db database.IDB,
	req *projectenvdto.DeleteProjectEnvReq,
	data *deleteProjectEnvData,
) error {
	// Loads and checks target project
	project, err := uc.projectRepo.GetByID(ctx, db, req.ProjectID,
		bunex.SelectFor("UPDATE OF project"),
		bunex.SelectRelation("Envs", bunex.SelectOrder("display_order")),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Project = project

	for _, env := range project.Envs {
		if env.DeletedAt.IsZero() && env.ID == req.ProjectEnvID {
			data.DeletingProjectEnv = env
		} else {
			data.UpdatingOrderEnvs = append(data.UpdatingOrderEnvs, env)
		}
	}
	if data.DeletingProjectEnv == nil {
		return apperrors.NewNotFound("ProjectEnv").
			WithMsgLog("project env %s not found in project %s", req.ProjectEnvID, req.ProjectID)
	}

	return nil
}

func (uc *ProjectEnvUC) prepareDeletingProjectEnv(
	envData *deleteProjectEnvData,
	persistingData *projectservice.PersistingProjectData,
) {
	timeNow := timeutil.NowUTC()

	// Deletes the env
	envData.DeletingProjectEnv.DeletedAt = timeNow
	persistingData.UpsertingEnvs = append(persistingData.UpsertingEnvs, envData.DeletingProjectEnv)

	// Updates order of the active envs
	for i, projectEnv := range envData.UpdatingOrderEnvs {
		if projectEnv.DisplayOrder != i {
			projectEnv.DisplayOrder = i
			persistingData.UpsertingEnvs = append(persistingData.UpsertingEnvs, projectEnv)
		}
	}
}
