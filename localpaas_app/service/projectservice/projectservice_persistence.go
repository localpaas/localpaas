package projectservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type PersistingProjectData struct {
	UpsertingProjects []*entity.Project
	UpsertingApps     []*entity.App
	UpsertingTags     []*entity.ProjectTag
	UpsertingSettings []*entity.Setting
	UpsertingAccesses []*entity.ACLPermission

	ProjectsToDeleteTags []string
}

func (s *projectService) PersistProjectData(ctx context.Context, db database.IDB,
	persistingData *PersistingProjectData) error {
	// Deletes all current linked data if configured
	err := s.projectTagRepo.DeleteAllByProjects(ctx, db, persistingData.ProjectsToDeleteTags)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Persists data
	// Settings
	err = s.settingRepo.UpsertMulti(ctx, db, persistingData.UpsertingSettings,
		entity.SettingUpsertingConflictCols, entity.SettingUpsertingUpdateCols)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Projects
	err = s.projectRepo.UpsertMulti(ctx, db, persistingData.UpsertingProjects,
		entity.ProjectUpsertingConflictCols, entity.ProjectUpsertingUpdateCols)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Apps
	err = s.appRepo.UpsertMulti(ctx, db, persistingData.UpsertingApps,
		entity.AppUpsertingConflictCols, entity.AppUpsertingUpdateCols)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Project Tags
	err = s.projectTagRepo.UpsertMulti(ctx, db, persistingData.UpsertingTags,
		entity.ProjectTagUpsertingConflictCols, entity.ProjectTagUpsertingUpdateCols)
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
