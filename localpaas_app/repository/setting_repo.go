package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/tiendc/gofn"
	"github.com/uptrace/bun"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
)

type SettingRepo interface {
	// Get a setting by id in the given scope with handling inheritance
	GetByID(ctx context.Context, db database.IDB, typ base.SettingType, id string, requireActive bool,
		opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetByIDGlobally(ctx context.Context, db database.IDB, typ base.SettingType, id string,
		requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetByIDAndProject(ctx context.Context, db database.IDB, typ base.SettingType, id, projectID string,
		requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetByIDAndApp(ctx context.Context, db database.IDB, typ base.SettingType, id, appID, projectID string,
		requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetByIDAndAppObject(ctx context.Context, db database.IDB, typ base.SettingType, id string, app *entity.App,
		requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetByIDAndUser(ctx context.Context, db database.IDB, typ base.SettingType, id, userID string,
		requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error)

	GetByKind(ctx context.Context, db database.IDB, typ base.SettingType, kind string, requireActive bool,
		opts ...bunex.SelectQueryOption) (*entity.Setting, error)

	// Get a setting by name in the given scope with handling inheritance
	GetByName(ctx context.Context, db database.IDB, typ base.SettingType, name string, requireActive bool,
		opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetByNameGlobally(ctx context.Context, db database.IDB, typ base.SettingType, name string, requireActive bool,
		opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetByNameAndProject(ctx context.Context, db database.IDB, typ base.SettingType, name, projectID string,
		requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetByNameAndApp(ctx context.Context, db database.IDB, typ base.SettingType, name, appID, projectID string,
		requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetByNameAndUser(ctx context.Context, db database.IDB, typ base.SettingType, name, userID string,
		requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error)

	// Get a single setting in the given scope with handling inheritance by priority order:
	// own setting, parent app setting, parent project setting, global setting.
	GetSingleGlobally(ctx context.Context, db database.IDB, typ base.SettingType,
		requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetSingleByProject(ctx context.Context, db database.IDB, typ base.SettingType, projectID string,
		requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetSingleByApp(ctx context.Context, db database.IDB, typ base.SettingType, appID, projectID string,
		requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetSingleByAppObject(ctx context.Context, db database.IDB, typ base.SettingType, app *entity.App,
		requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetSingleByUser(ctx context.Context, db database.IDB, typ base.SettingType, userID string,
		requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error)

	// List settings in the given scope with handling inheritance
	List(ctx context.Context, db database.IDB, paging *basedto.Paging,
		opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error)
	ListGlobally(ctx context.Context, db database.IDB, paging *basedto.Paging,
		opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error)
	ListByProject(ctx context.Context, db database.IDB, projectID string, paging *basedto.Paging,
		opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error)
	ListByApp(ctx context.Context, db database.IDB, appID, projectID string, paging *basedto.Paging,
		opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error)
	ListByAppObject(ctx context.Context, db database.IDB, app *entity.App, paging *basedto.Paging,
		opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error)
	ListByUser(ctx context.Context, db database.IDB, userID string, paging *basedto.Paging,
		opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error)

	ListByIDs(ctx context.Context, db database.IDB, ids []string, requireActive bool,
		opts ...bunex.SelectQueryOption) ([]*entity.Setting, error)
	ListByIDsAsMap(ctx context.Context, db database.IDB, ids []string, requireActive bool,
		opts ...bunex.SelectQueryOption) (map[string]*entity.Setting, error)

	// Make sure there is at most one active setting in the given scope.
	// Inherited settings are still allowed.
	EnsureSingleGlobally(ctx context.Context, db database.IDB, typ base.SettingType,
		opts ...bunex.SelectQueryOption) error
	EnsureSingleByProject(ctx context.Context, db database.IDB, typ base.SettingType, projectID string,
		opts ...bunex.SelectQueryOption) error
	EnsureSingleByApp(ctx context.Context, db database.IDB, typ base.SettingType, appID string,
		opts ...bunex.SelectQueryOption) error
	EnsureSingleByUser(ctx context.Context, db database.IDB, typ base.SettingType, userID string,
		opts ...bunex.SelectQueryOption) error

	Insert(ctx context.Context, db database.IDB, setting *entity.Setting,
		opts ...bunex.InsertQueryOption) error
	Upsert(ctx context.Context, db database.IDB, setting *entity.Setting,
		conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error
	UpsertMulti(ctx context.Context, db database.IDB, settings []*entity.Setting,
		conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error
	Update(ctx context.Context, db database.IDB, setting *entity.Setting,
		opts ...bunex.UpdateQueryOption) error
	UpdateClearDefaultFlag(ctx context.Context, db database.IDB, typ base.SettingType, exceptID string,
		opts ...bunex.UpdateQueryOption) error
}

type settingRepo struct {
	appRepo AppRepo
}

func NewSettingRepo(appRepo AppRepo) SettingRepo {
	return &settingRepo{
		appRepo: appRepo,
	}
}

func (repo *settingRepo) GetByID(ctx context.Context, db database.IDB, typ base.SettingType, id string,
	requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	opts = repo.applyFilter(opts, typ, id, requireActive)
	return repo.get(ctx, db, opts...)
}

func (repo *settingRepo) GetByIDGlobally(ctx context.Context, db database.IDB, typ base.SettingType, id string,
	requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	opts = repo.applyFilter(opts, typ, id, requireActive)
	opts = append(opts, bunex.SelectWhere("setting.object_id IS NULL"))
	return repo.get(ctx, db, opts...)
}

func (repo *settingRepo) GetByIDAndProject(ctx context.Context, db database.IDB, typ base.SettingType,
	id, projectID string, requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	opts = repo.applyFilter(opts, typ, id, requireActive)
	opts = repo.applyAppAndProjectFilter(opts, "", "", projectID)
	return repo.get(ctx, db, opts...)
}

func (repo *settingRepo) GetByIDAndApp(ctx context.Context, db database.IDB, typ base.SettingType,
	id, appID, projectID string, requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	var parentAppID string
	if appID != "" {
		app, err := repo.appRepo.GetByID(ctx, db, projectID, appID,
			bunex.SelectColumns("project_id", "parent_id"))
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		parentAppID = app.ParentID
		projectID = app.ProjectID
	}
	opts = repo.applyFilter(opts, typ, id, requireActive)
	opts = repo.applyAppAndProjectFilter(opts, appID, parentAppID, projectID)
	return repo.get(ctx, db, opts...)
}

func (repo *settingRepo) GetByIDAndAppObject(ctx context.Context, db database.IDB, typ base.SettingType,
	id string, app *entity.App, requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	opts = repo.applyFilter(opts, typ, id, requireActive)
	opts = repo.applyAppAndProjectFilter(opts, app.ID, app.ParentID, app.ProjectID)
	return repo.get(ctx, db, opts...)
}

func (repo *settingRepo) GetByIDAndUser(ctx context.Context, db database.IDB, typ base.SettingType,
	id, userID string, requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	opts = repo.applyFilter(opts, typ, id, requireActive)
	opts = repo.applyUserFilter(opts, userID)
	return repo.get(ctx, db, opts...)
}

func (repo *settingRepo) GetByKind(ctx context.Context, db database.IDB, typ base.SettingType, kind string,
	requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	if kind == "" {
		return nil, nil
	}
	opts = repo.applyFilter(opts, typ, "", requireActive)
	opts = repo.applyNameAndKindFilter(opts, "", kind)
	return repo.get(ctx, db, opts...)
}

func (repo *settingRepo) GetByName(ctx context.Context, db database.IDB, typ base.SettingType, name string,
	requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	if name == "" {
		return nil, nil
	}
	opts = repo.applyFilter(opts, typ, "", requireActive)
	opts = repo.applyNameAndKindFilter(opts, name, "")
	return repo.get(ctx, db, opts...)
}

func (repo *settingRepo) GetByNameGlobally(ctx context.Context, db database.IDB, typ base.SettingType, name string,
	requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	if name == "" {
		return nil, nil
	}
	opts = repo.applyFilter(opts, typ, "", requireActive)
	opts = repo.applyNameAndKindFilter(opts, name, "")
	opts = append(opts, bunex.SelectWhere("setting.object_id IS NULL"))
	return repo.get(ctx, db, opts...)
}

func (repo *settingRepo) GetByNameAndProject(ctx context.Context, db database.IDB, typ base.SettingType,
	name, projectID string, requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	if name == "" {
		return nil, nil
	}
	opts = repo.applyFilter(opts, typ, "", requireActive)
	opts = repo.applyNameAndKindFilter(opts, name, "")
	opts = repo.applyAppAndProjectFilter(opts, "", "", projectID)
	return repo.get(ctx, db, opts...)
}

func (repo *settingRepo) GetByNameAndApp(ctx context.Context, db database.IDB, typ base.SettingType,
	name, appID, projectID string, requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	if name == "" {
		return nil, nil
	}
	var parentAppID string
	if appID != "" {
		app, err := repo.appRepo.GetByID(ctx, db, projectID, appID,
			bunex.SelectColumns("project_id", "parent_id"))
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		parentAppID = app.ParentID
		projectID = app.ProjectID
	}
	opts = repo.applyFilter(opts, typ, "", requireActive)
	opts = repo.applyNameAndKindFilter(opts, name, "")
	opts = repo.applyAppAndProjectFilter(opts, appID, parentAppID, projectID)
	return repo.get(ctx, db, opts...)
}

func (repo *settingRepo) GetByNameAndAppObject(ctx context.Context, db database.IDB, typ base.SettingType,
	name string, app *entity.App, requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	if name == "" {
		return nil, nil
	}
	opts = repo.applyFilter(opts, typ, "", requireActive)
	opts = repo.applyNameAndKindFilter(opts, name, "")
	opts = repo.applyAppAndProjectFilter(opts, app.ID, app.ParentID, app.ProjectID)
	return repo.get(ctx, db, opts...)
}

func (repo *settingRepo) GetByNameAndUser(ctx context.Context, db database.IDB, typ base.SettingType,
	name, userID string, requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	if name == "" {
		return nil, nil
	}
	opts = repo.applyFilter(opts, typ, "", requireActive)
	opts = repo.applyNameAndKindFilter(opts, name, "")
	opts = repo.applyUserFilter(opts, userID)
	return repo.get(ctx, db, opts...)
}

func (repo *settingRepo) get(ctx context.Context, db database.IDB,
	opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	setting := &entity.Setting{}
	query := db.NewSelect().Model(setting)
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if setting == nil || errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.NewNotFound("Setting").WithCause(err)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	if hasChange, _ := repo.updateExpiredSetting(ctx, db, setting); hasChange {
		return repo.get(ctx, db, opts...)
	}
	return setting, nil
}

func (repo *settingRepo) GetSingleGlobally(ctx context.Context, db database.IDB, typ base.SettingType,
	requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	opts = repo.applyFilter(opts, typ, "", requireActive)
	setting, err := repo.get(ctx, db, opts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return setting, nil
}

func (repo *settingRepo) GetSingleByProject(ctx context.Context, db database.IDB, typ base.SettingType,
	projectID string, requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	opts = repo.applyFilter(opts, typ, "", requireActive)
	settings, _, err := repo.ListByProject(ctx, db, projectID, nil, opts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if len(settings) == 0 {
		return nil, apperrors.NewNotFound("Setting")
	}
	if len(settings) == 1 {
		return settings[0], nil
	}

	var globalSetting *entity.Setting
	for _, setting := range settings {
		if setting.ObjectID == projectID { // setting belongs to the project directly has the highest priority
			return setting, nil
		}
		if globalSetting == nil && setting.ObjectID == "" {
			globalSetting = setting
			continue
		}
	}
	if globalSetting != nil {
		return globalSetting, nil
	}
	return nil, apperrors.NewNotFound("Setting")
}

func (repo *settingRepo) GetSingleByApp(ctx context.Context, db database.IDB, typ base.SettingType,
	appID, projectID string, requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	app := &entity.App{
		ID:        appID,
		ProjectID: projectID,
	}
	if appID != "" {
		currApp, err := repo.appRepo.GetByID(ctx, db, projectID, appID,
			bunex.SelectColumns("project_id", "parent_id"))
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		app.ParentID = currApp.ParentID
		app.ProjectID = currApp.ProjectID
	}
	return repo.GetSingleByAppObject(ctx, db, typ, app, requireActive, opts...)
}

func (repo *settingRepo) GetSingleByAppObject(ctx context.Context, db database.IDB, typ base.SettingType,
	app *entity.App, requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	opts = repo.applyFilter(opts, typ, "", requireActive)
	settings, _, err := repo.ListByAppObject(ctx, db, app, nil, opts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if len(settings) == 0 {
		return nil, apperrors.NewNotFound("Setting")
	}
	if len(settings) == 1 {
		return settings[0], nil
	}

	var parentSetting, projectSetting, globalSetting *entity.Setting
	for _, setting := range settings {
		if setting.ObjectID == app.ID { // setting belongs to the app directly has the highest priority
			return setting, nil
		}
		if parentSetting == nil && setting.ObjectID == app.ParentID {
			parentSetting = setting
			continue
		}
		if projectSetting == nil && setting.ObjectID == app.ProjectID {
			projectSetting = setting
			continue
		}
		if globalSetting == nil && setting.ObjectID == "" {
			globalSetting = setting
			continue
		}
	}
	setting := gofn.Coalesce(parentSetting, projectSetting, globalSetting)
	if setting != nil {
		return setting, nil
	}
	return nil, apperrors.NewNotFound("Setting")
}

func (repo *settingRepo) GetSingleByUser(ctx context.Context, db database.IDB, typ base.SettingType,
	userID string, requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	opts = repo.applyFilter(opts, typ, "", requireActive)
	settings, _, err := repo.ListByUser(ctx, db, userID, nil, opts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if len(settings) == 0 {
		return nil, apperrors.NewNotFound("Setting")
	}
	if len(settings) == 1 {
		return settings[0], nil
	}

	var globalSetting *entity.Setting
	for _, setting := range settings {
		if setting.ObjectID == userID { // setting belongs to the user directly has the highest priority
			return setting, nil
		}
		if globalSetting == nil && setting.ObjectID == "" {
			globalSetting = setting
			continue
		}
	}
	if globalSetting != nil {
		return globalSetting, nil
	}
	return nil, apperrors.NewNotFound("Setting")
}

func (repo *settingRepo) List(ctx context.Context, db database.IDB, paging *basedto.Paging,
	opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error) {
	var settings []*entity.Setting
	query := db.NewSelect().Model(&settings)
	query = bunex.ApplySelect(query, opts...)

	var pagingMeta *basedto.PagingMeta
	if paging != nil {
		pagingMeta = newPagingMeta(paging)

		// Counts the total first
		total, err := query.Count(ctx)
		if err != nil {
			return nil, nil, apperrors.Wrap(err)
		}
		pagingMeta.Total = total

		// Applies pagination
		query = bunex.ApplyPagination(query, paging)
	}

	err := query.Scan(ctx)
	if err != nil {
		return nil, nil, wrapPaginationError(err, paging)
	}

	if hasChange, _ := repo.updateExpiredSettings(ctx, db, settings); hasChange {
		return repo.List(ctx, db, paging, opts...)
	}
	return settings, pagingMeta, nil
}

func (repo *settingRepo) ListGlobally(ctx context.Context, db database.IDB,
	paging *basedto.Paging, opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error) {
	opts = append(opts, bunex.SelectWhere("setting.object_id IS NULL"))
	return repo.List(ctx, db, paging, opts...)
}

func (repo *settingRepo) ListByProject(ctx context.Context, db database.IDB, projectID string,
	paging *basedto.Paging, opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error) {
	opts = repo.applyAppAndProjectFilter(opts, "", "", projectID)
	return repo.List(ctx, db, paging, opts...)
}

func (repo *settingRepo) ListByApp(ctx context.Context, db database.IDB, appID, projectID string,
	paging *basedto.Paging, opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error) {
	var parentAppID string
	if appID != "" {
		// Query app to get its parent ID if there is
		app, err := repo.appRepo.GetByID(ctx, db, projectID, appID,
			bunex.SelectColumns("project_id", "parent_id"))
		if err != nil {
			return nil, nil, apperrors.Wrap(err)
		}
		parentAppID = app.ParentID
		projectID = app.ProjectID
	}
	opts = repo.applyAppAndProjectFilter(opts, appID, parentAppID, projectID)
	return repo.List(ctx, db, paging, opts...)
}

func (repo *settingRepo) ListByAppObject(ctx context.Context, db database.IDB, app *entity.App,
	paging *basedto.Paging, opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error) {
	opts = repo.applyAppAndProjectFilter(opts, app.ID, app.ParentID, app.ProjectID)
	return repo.List(ctx, db, paging, opts...)
}

func (repo *settingRepo) ListByUser(ctx context.Context, db database.IDB, userID string,
	paging *basedto.Paging, opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error) {
	opts = repo.applyUserFilter(opts, userID)
	return repo.List(ctx, db, paging, opts...)
}

func (repo *settingRepo) ListByIDs(ctx context.Context, db database.IDB, ids []string, requireActive bool,
	opts ...bunex.SelectQueryOption) ([]*entity.Setting, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	opts = append(opts, bunex.SelectWhere("setting.id IN (?)", bun.In(ids)))
	opts = repo.applyFilter(opts, "", "", requireActive)
	settings, _, err := repo.List(ctx, db, nil, opts...)
	return settings, apperrors.Wrap(err)
}

func (repo *settingRepo) ListByIDsAsMap(ctx context.Context, db database.IDB, ids []string, requireActive bool,
	opts ...bunex.SelectQueryOption) (map[string]*entity.Setting, error) {
	settings, err := repo.ListByIDs(ctx, db, ids, requireActive, opts...)
	if err != nil {
		return nil, apperrors.New(err)
	}
	res := make(map[string]*entity.Setting, len(settings))
	for _, setting := range settings {
		res[setting.ID] = setting
	}
	return res, nil
}

func (repo *settingRepo) applyFilter(opts []bunex.SelectQueryOption, typ base.SettingType, id string,
	requireActive bool) []bunex.SelectQueryOption {
	if typ != "" {
		opts = append(opts, bunex.SelectWhere("setting.type = ?", typ))
	}
	if id != "" {
		opts = append(opts, bunex.SelectWhere("setting.id = ?", id))
	}
	if requireActive {
		opts = append(opts, bunex.SelectWhere("setting.status = ?", base.SettingStatusActive))
	}
	return opts
}

func (repo *settingRepo) applyNameAndKindFilter(opts []bunex.SelectQueryOption,
	name, kind string) []bunex.SelectQueryOption {
	if name != "" {
		opts = append(opts, bunex.SelectWhere("LOWER(setting.name) = ?", strings.ToLower(name)))
	}
	if kind != "" {
		opts = append(opts, bunex.SelectWhere("setting.kind = ?", kind))
	}
	return opts
}

func (repo *settingRepo) applyAppAndProjectFilter(opts []bunex.SelectQueryOption,
	appID, parentAppID, projectID string) []bunex.SelectQueryOption {
	if projectID != "" {
		opts = append(opts,
			bunex.SelectJoin("LEFT JOIN project_shared_settings pss ON pss.setting_id = setting.id"),
			bunex.SelectWhereGroup(
				bunex.SelectWhereIf(appID != "", "setting.object_id = ?", appID),
				bunex.SelectWhereOrIf(parentAppID != "", "setting.object_id = ?", parentAppID),
				bunex.SelectWhereOr("setting.object_id = ?", projectID),
				bunex.SelectWhereOr("(setting.object_id IS NULL AND setting.avail_in_projects = TRUE)"),
				bunex.SelectWhereOr("(setting.object_id IS NULL AND pss.project_id = ?)", projectID),
			),
		)
	} else if appID != "" {
		opts = append(opts,
			bunex.SelectWhereGroup(
				bunex.SelectWhere("setting.object_id = ?", appID),
				bunex.SelectWhereOrIf(parentAppID != "", "setting.object_id = ?", parentAppID),
			),
		)
	}
	return opts
}

func (repo *settingRepo) applyUserFilter(opts []bunex.SelectQueryOption, userID string) []bunex.SelectQueryOption {
	opts = append(opts, bunex.SelectWhere("setting.object_id = ?", userID))
	return opts
}

func (repo *settingRepo) EnsureSingleGlobally(ctx context.Context, db database.IDB, typ base.SettingType,
	opts ...bunex.SelectQueryOption) error {
	query := db.NewSelect().Model((*entity.Setting)(nil)).
		Where("setting.type = ?", typ).
		Where("setting.status = ?", base.SettingStatusActive).
		Where("setting.object_id IS NULL")
	query = bunex.ApplySelect(query, opts...)

	count, err := query.Count(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return apperrors.Wrap(err)
	}
	if count > 1 {
		return apperrors.Wrap(apperrors.ErrConflict)
	}
	return nil
}

func (repo *settingRepo) EnsureSingleByProject(ctx context.Context, db database.IDB, typ base.SettingType,
	projectID string, opts ...bunex.SelectQueryOption) error {
	query := db.NewSelect().Model((*entity.Setting)(nil)).
		Where("setting.type = ?", typ).
		Where("setting.status = ?", base.SettingStatusActive).
		Where("setting.object_id = ?", projectID)
	query = bunex.ApplySelect(query, opts...)

	count, err := query.Count(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return apperrors.Wrap(err)
	}
	if count > 1 {
		return apperrors.Wrap(apperrors.ErrConflict)
	}
	return nil
}

func (repo *settingRepo) EnsureSingleByApp(ctx context.Context, db database.IDB, typ base.SettingType,
	appID string, opts ...bunex.SelectQueryOption) error {
	query := db.NewSelect().Model((*entity.Setting)(nil)).
		Where("setting.type = ?", typ).
		Where("setting.status = ?", base.SettingStatusActive).
		Where("setting.object_id = ?", appID)
	query = bunex.ApplySelect(query, opts...)

	count, err := query.Count(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return apperrors.Wrap(err)
	}
	if count > 1 {
		return apperrors.Wrap(apperrors.ErrConflict)
	}
	return nil
}

func (repo *settingRepo) EnsureSingleByUser(ctx context.Context, db database.IDB, typ base.SettingType,
	userID string, opts ...bunex.SelectQueryOption) error {
	query := db.NewSelect().Model((*entity.Setting)(nil)).
		Where("setting.type = ?", typ).
		Where("setting.status = ?", base.SettingStatusActive).
		Where("setting.object_id = ?", userID)
	query = bunex.ApplySelect(query, opts...)

	count, err := query.Count(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return apperrors.Wrap(err)
	}
	if count > 1 {
		return apperrors.Wrap(apperrors.ErrConflict)
	}
	return nil
}

func (repo *settingRepo) Insert(ctx context.Context, db database.IDB, setting *entity.Setting,
	opts ...bunex.InsertQueryOption) error {
	query := db.NewInsert().Model(setting)
	query = bunex.ApplyInsert(query, opts...)

	_, err := query.Exec(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (repo *settingRepo) Upsert(ctx context.Context, db database.IDB, setting *entity.Setting,
	conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error {
	return repo.UpsertMulti(ctx, db, []*entity.Setting{setting}, conflictCols, updateCols, opts...)
}

func (repo *settingRepo) UpsertMulti(ctx context.Context, db database.IDB, settings []*entity.Setting,
	conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error {
	if len(settings) == 0 {
		return nil
	}
	query := db.NewInsert().Model(&settings)
	query = bunex.ApplyInsert(query, opts...)
	query = bunex.ApplyUpsert(query, conflictCols, updateCols)

	_, err := query.Exec(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (repo *settingRepo) Update(ctx context.Context, db database.IDB, setting *entity.Setting,
	opts ...bunex.UpdateQueryOption) error {
	query := db.NewUpdate().Model(setting).WherePK()
	query = bunex.ApplyUpdate(query, opts...)

	_, err := query.Exec(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (repo *settingRepo) UpdateClearDefaultFlag(ctx context.Context, db database.IDB, typ base.SettingType,
	exceptID string, opts ...bunex.UpdateQueryOption) error {
	query := db.NewUpdate().Model((*entity.Setting)(nil)).
		Where("setting.type = ?", typ).
		Where("setting.is_default = true").
		Set("is_default = false")
	if exceptID != "" {
		query = query.Where("setting.id != ?", exceptID)
	}
	query = bunex.ApplyUpdate(query, opts...)

	_, err := query.Exec(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (repo *settingRepo) updateExpiredSetting(ctx context.Context, db database.IDB,
	setting *entity.Setting) (bool, error) {
	if setting == nil {
		return false, nil
	}
	return repo.updateExpiredSettings(ctx, db, []*entity.Setting{setting})
}

func (repo *settingRepo) updateExpiredSettings(ctx context.Context, db database.IDB,
	settings []*entity.Setting) (hasChange bool, err error) {
	for _, setting := range settings {
		if setting.IsStatusDirty() {
			hasChange = true
			break
		}
	}
	if !hasChange {
		return false, nil
	}
	query := db.NewUpdate().Model((*entity.Setting)(nil)).
		Set("status = ?", base.SettingStatusExpired).
		Where("status = ? AND expire_at < NOW()", base.SettingStatusActive)

	_, err = query.Exec(ctx)
	if err != nil {
		return hasChange, apperrors.Wrap(err)
	}
	return hasChange, nil
}
