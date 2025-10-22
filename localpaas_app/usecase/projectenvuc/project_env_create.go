package projectenvuc

import (
	"context"
	"strings"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
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

func (uc *ProjectEnvUC) CreateProjectEnv(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectenvdto.CreateProjectEnvReq,
) (*projectenvdto.CreateProjectEnvResp, error) {
	var persistingData *projectservice.PersistingProjectData
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		envData := &createProjectEnvData{}
		err := uc.loadProjectEnvData(ctx, db, req, envData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData = &projectservice.PersistingProjectData{}
		uc.preparePersistingProjectEnvs(auth, envData.Project, []*projectenvdto.ProjectEnvReq{req.ProjectEnvReq},
			envData.NextDisplayOrder, timeutil.NowUTC(), persistingData)

		return uc.projectService.PersistProjectData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &projectenvdto.CreateProjectEnvResp{
		Data: &basedto.ObjectIDResp{ID: persistingData.UpsertingEnvs[0].ID},
	}, nil
}

type createProjectEnvData struct {
	Project          *entity.Project
	NextDisplayOrder int
}

func (uc *ProjectEnvUC) loadProjectEnvData(
	ctx context.Context,
	db database.IDB,
	req *projectenvdto.CreateProjectEnvReq,
	data *createProjectEnvData,
) error {
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

func (uc *ProjectEnvUC) preparePersistingProjectEnvs(
	auth *basedto.Auth,
	project *entity.Project,
	envReqs []*projectenvdto.ProjectEnvReq,
	startDisplayOrder int,
	timeNow time.Time,
	persistingData *projectservice.PersistingProjectData,
) {
	displayOrder := startDisplayOrder
	for _, envReq := range envReqs {
		persistingData.UpsertingEnvs = append(persistingData.UpsertingEnvs,
			&entity.ProjectEnv{
				ID:           gofn.Must(ulid.NewStringULID()),
				ProjectID:    project.ID,
				Name:         envReq.Name,
				Status:       envReq.Status,
				DisplayOrder: displayOrder,
				CreatedAt:    timeNow,
				CreatedBy:    auth.User.ID,
				UpdatedAt:    timeNow,
				UpdatedBy:    auth.User.ID,
			})
		displayOrder++
	}
}
