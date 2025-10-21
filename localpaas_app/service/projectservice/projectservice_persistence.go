package projectservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type PersistingProjectData struct {
	UpsertingProjects []*entity.Project
	UpsertingTags     []*entity.ProjectTag
	UpsertingEnvs     []*entity.ProjectEnv
	UpsertingSettings []*entity.Setting
	UpsertingAccesses []*entity.ACLPermission

	ProjectsToDeleteTags     []string
	ProjectsToDeleteEnvs     []string
	ProjectsToDeleteSettings []string
}

func (s *projectService) PersistProjectData(ctx context.Context, db database.IDB,
	persistingData *PersistingProjectData) error {
	// Deletes all current linked data if configured
	err := s.projectTagRepo.DeleteAllByProjects(ctx, db, persistingData.ProjectsToDeleteTags)
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = s.projectEnvRepo.DeleteAllByProjects(ctx, db, persistingData.ProjectsToDeleteEnvs)
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = s.settingRepo.DeleteAllByTargetObjects(ctx, db, base.SettingTargetProject,
		persistingData.ProjectsToDeleteSettings)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Persists data
	err = s.projectRepo.UpsertMulti(ctx, db, persistingData.UpsertingProjects,
		entity.ProjectUpsertingConflictCols, entity.ProjectUpsertingUpdateCols)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Tags
	err = s.projectTagRepo.UpsertMulti(ctx, db, persistingData.UpsertingTags,
		entity.ProjectTagUpsertingConflictCols, entity.ProjectTagUpsertingUpdateCols)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Envs
	err = s.projectEnvRepo.UpsertMulti(ctx, db, persistingData.UpsertingEnvs,
		entity.ProjectEnvUpsertingConflictCols, entity.ProjectEnvUpsertingUpdateCols)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Project accesses
	err = s.permissionManager.UpdateACLPermissions(ctx, db, persistingData.UpsertingAccesses)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
