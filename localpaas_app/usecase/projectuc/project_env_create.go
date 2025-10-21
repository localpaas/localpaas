package projectuc

import (
	"context"
	"strings"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc/projectdto"
	"github.com/localpaas/localpaas/pkg/timeutil"
)

func (uc *ProjectUC) CreateProjectEnv(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectdto.CreateProjectEnvReq,
) (*projectdto.CreateProjectEnvResp, error) {
	var persistingData *persistingProjectData
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		envData := &createProjectEnvData{}
		err := uc.loadProjectEnvDataForAddNew(ctx, db, req, envData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData = &persistingProjectData{}
		uc.preparePersistingProjectEnvs(auth, envData.Project, []string{req.Name}, envData.NextDisplayOrder,
			timeutil.NowUTC(), persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &projectdto.CreateProjectEnvResp{
		Data: &basedto.ObjectIDResp{ID: persistingData.UpsertingEnvs[0].ID},
	}, nil
}

type createProjectEnvData struct {
	Project          *entity.Project
	NextDisplayOrder int
}

func (uc *ProjectUC) loadProjectEnvDataForAddNew(
	ctx context.Context,
	db database.IDB,
	req *projectdto.CreateProjectEnvReq,
	data *createProjectEnvData,
) error {
	// Loads and checks target project
	project, err := uc.projectRepo.GetByID(ctx, db, req.ProjectID,
		bunex.SelectFor("UPDATE OF project"),
		bunex.SelectRelation("Envs"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Project = project

	nextDisplayOrder := 0
	for _, projectEnv := range project.Envs {
		if projectEnv.DeletedAt.IsZero() && strings.EqualFold(projectEnv.Name, req.Name) {
			return apperrors.NewAlreadyExist("ProjectEnv")
		}
		nextDisplayOrder = max(nextDisplayOrder, projectEnv.DisplayOrder+1)
	}
	data.NextDisplayOrder = nextDisplayOrder

	return nil
}
