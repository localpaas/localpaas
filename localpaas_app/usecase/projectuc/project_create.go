package projectuc

import (
	"context"
	"errors"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/slugify"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/service/projectservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc/projectdto"
)

const (
	projectKeyMaxLen = 100
)

var (
	unallowedProjectKey = []string{"localpaas"}
)

func (uc *ProjectUC) CreateProject(
	ctx context.Context,
	auth *basedto.Auth,
	req *projectdto.CreateProjectReq,
) (*projectdto.CreateProjectResp, error) {
	projectData := &createProjectData{}
	err := uc.loadProjectData(ctx, uc.db, req, projectData)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	persistingData := &persistingProjectData{}
	uc.preparePersistingProject(auth, req, timeutil.NowUTC(), projectData, persistingData)

	err = transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	createdProject := persistingData.UpsertingProjects[0]

	// Create default network for the project
	_, err = uc.networkService.CreateProjectNetwork(ctx, createdProject)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &projectdto.CreateProjectResp{
		Data: &basedto.ObjectIDResp{ID: createdProject.ID},
	}, nil
}

type createProjectData struct {
	ProjectKey string
}

func (uc *ProjectUC) loadProjectData(
	ctx context.Context,
	db database.IDB,
	req *projectdto.CreateProjectReq,
	data *createProjectData,
) error {
	data.ProjectKey = slugify.SlugifyEx(req.Name, nil, projectKeyMaxLen)
	if gofn.Contain(unallowedProjectKey, data.ProjectKey) {
		return apperrors.New(apperrors.ErrNameUnavailable).WithMsgLog("project name is not allowed")
	}

	project, err := uc.projectRepo.GetByKey(ctx, db, data.ProjectKey)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if project != nil {
		return apperrors.NewAlreadyExist("Project").
			WithMsgLog("project key '%s' already exists", data.ProjectKey)
	}

	return nil
}

type persistingProjectData struct {
	projectservice.PersistingProjectData
}

func (uc *ProjectUC) preparePersistingProject(
	auth *basedto.Auth,
	req *projectdto.CreateProjectReq,
	timeNow time.Time,
	data *createProjectData,
	persistingData *persistingProjectData,
) {
	// Upserting project
	project := &entity.Project{
		ID:        gofn.Must(ulid.NewStringULID()),
		Key:       data.ProjectKey,
		CreatedAt: timeNow,
	}

	uc.preparePersistingProjectBase(auth, project, req.ProjectBaseReq, timeNow, persistingData)
	uc.preparePersistingProjectTags(project, req.Tags, 0, persistingData)
}

func (uc *ProjectUC) preparePersistingProjectBase(
	auth *basedto.Auth,
	project *entity.Project,
	req *projectdto.ProjectBaseReq,
	timeNow time.Time,
	persistingData *persistingProjectData,
) {
	project.Name = req.Name
	project.Status = req.Status
	project.Note = req.Note
	project.OwnerID = gofn.Coalesce(req.Owner.ID, auth.User.ID)
	project.UpdatedAt = timeNow

	persistingData.UpsertingProjects = append(persistingData.UpsertingProjects, project)
}

func (uc *ProjectUC) preparePersistingProjectTags(
	project *entity.Project,
	tags []string,
	startDisplayOrder int,
	persistingData *persistingProjectData,
) {
	displayOrder := startDisplayOrder
	for _, tag := range tags {
		persistingData.UpsertingTags = append(persistingData.UpsertingTags,
			&entity.ProjectTag{
				ProjectID:    project.ID,
				Tag:          tag,
				DisplayOrder: displayOrder,
			})
		displayOrder++
	}
}

func (uc *ProjectUC) persistData(
	ctx context.Context,
	db database.IDB,
	persistingData *persistingProjectData,
) error {
	err := uc.projectService.PersistProjectData(ctx, db, &persistingData.PersistingProjectData)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
