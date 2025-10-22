package projectuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc/projectdto"
)

func (uc *ProjectUC) DeleteProject(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectdto.DeleteProjectReq,
) (*projectdto.DeleteProjectResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		projectData := &deleteProjectData{}
		err := uc.loadProjectDataForDelete(ctx, db, req, projectData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingProjectData{}
		uc.prepareDeletingProject(projectData, persistingData)

		// TODO: handle project deletion

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &projectdto.DeleteProjectResp{}, nil
}

type deleteProjectData struct {
	Project *entity.Project
}

func (uc *ProjectUC) loadProjectDataForDelete(
	ctx context.Context,
	db database.IDB,
	req *projectdto.DeleteProjectReq,
	data *deleteProjectData,
) error {
	project, err := uc.projectRepo.GetByID(ctx, db, req.ProjectID,
		bunex.SelectFor("UPDATE OF project"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Project = project

	if project.Status == base.ProjectStatusDeleting { //nolint
		// TODO: handle task deletion if previously failed
	}

	return nil
}

func (uc *ProjectUC) prepareDeletingProject(
	data *deleteProjectData,
	persistingData *persistingProjectData,
) {
	project := data.Project
	project.Status = base.ProjectStatusDeleting
	persistingData.UpsertingProjects = append(persistingData.UpsertingProjects, project)
}
