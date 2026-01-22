package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
)

type SettingRepo interface {
	GetByID(ctx context.Context, db database.IDB, typ base.SettingType, id string, requireActive bool,
		opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetByIDAndProject(ctx context.Context, db database.IDB, typ base.SettingType, id, projectID string,
		requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetByIDAndApp(ctx context.Context, db database.IDB, typ base.SettingType, id, projectID, appID string,
		requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetByIDAndUser(ctx context.Context, db database.IDB, typ base.SettingType, id, userID string,
		requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error)

	GetByKind(ctx context.Context, db database.IDB, typ base.SettingType, kind string, requireActive bool,
		opts ...bunex.SelectQueryOption) (*entity.Setting, error)

	GetByName(ctx context.Context, db database.IDB, typ base.SettingType, name string, requireActive bool,
		opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetByNameAndProject(ctx context.Context, db database.IDB, typ base.SettingType, name, projectID string,
		requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetByNameAndApp(ctx context.Context, db database.IDB, typ base.SettingType, name, projectID, appID string,
		requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetByNameAndUser(ctx context.Context, db database.IDB, typ base.SettingType, name, userID string,
		requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error)

	List(ctx context.Context, db database.IDB, paging *basedto.Paging,
		opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error)
	ListByProject(ctx context.Context, db database.IDB, projectID string, paging *basedto.Paging,
		opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error)
	ListByApp(ctx context.Context, db database.IDB, projectID, appID string, paging *basedto.Paging,
		opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error)
	ListByUser(ctx context.Context, db database.IDB, userID string, paging *basedto.Paging,
		opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error)

	ListByIDs(ctx context.Context, db database.IDB, ids []string, requireActive bool,
		opts ...bunex.SelectQueryOption) ([]*entity.Setting, error)
	ListByIDsAsMap(ctx context.Context, db database.IDB, ids []string, requireActive bool,
		opts ...bunex.SelectQueryOption) (map[string]*entity.Setting, error)

	Upsert(ctx context.Context, db database.IDB, setting *entity.Setting,
		conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error
	UpsertMulti(ctx context.Context, db database.IDB, settings []*entity.Setting,
		conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error
	Update(ctx context.Context, db database.IDB, setting *entity.Setting,
		opts ...bunex.UpdateQueryOption) error
}

type settingRepo struct {
}

func NewSettingRepo() SettingRepo {
	return &settingRepo{}
}

func (repo *settingRepo) GetByID(ctx context.Context, db database.IDB, typ base.SettingType, id string,
	requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	setting := &entity.Setting{}
	query := db.NewSelect().Model(setting).Where("setting.id = ?", id)
	if typ != "" {
		query = query.Where("setting.type = ?", typ)
	}
	if requireActive {
		query = query.Where("setting.status = ?", base.SettingStatusActive)
	}
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if setting == nil || errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.NewNotFound("Setting").WithCause(err)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	if hasChange, _ := repo.updateExpiredSetting(ctx, db, setting); hasChange {
		return repo.GetByID(ctx, db, typ, id, requireActive, opts...)
	}
	return setting, nil
}

func (repo *settingRepo) GetByIDAndProject(ctx context.Context, db database.IDB, typ base.SettingType,
	id, projectID string, requireActive bool, opts ...bunex.SelectQueryOption) (_ *entity.Setting, err error) {
	if projectID == "" {
		return repo.GetByID(ctx, db, typ, id, requireActive, opts...)
	}

	opts = append(opts, bunex.SelectWhere("setting.id = ?", id))
	if typ != "" {
		opts = append(opts, bunex.SelectWhere("setting.type = ?", typ))
	}
	if requireActive {
		opts = append(opts, bunex.SelectWhere("setting.status = ?", base.SettingStatusActive))
	}

	settings, _, err := repo.ListByProject(ctx, db, projectID, nil, opts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if len(settings) == 0 {
		return nil, apperrors.NewNotFound("Setting")
	}

	return settings[0], nil
}

func (repo *settingRepo) GetByIDAndApp(ctx context.Context, db database.IDB, typ base.SettingType,
	id, projectID, appID string, requireActive bool,
	opts ...bunex.SelectQueryOption) (_ *entity.Setting, err error) {
	if projectID == "" && appID == "" {
		return repo.GetByID(ctx, db, typ, id, requireActive, opts...)
	}

	opts = append(opts, bunex.SelectWhere("setting.id = ?", id))
	if typ != "" {
		opts = append(opts, bunex.SelectWhere("setting.type = ?", typ))
	}
	if requireActive {
		opts = append(opts, bunex.SelectWhere("setting.status = ?", base.SettingStatusActive))
	}

	settings, _, err := repo.ListByApp(ctx, db, projectID, appID, nil, opts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if len(settings) == 0 {
		return nil, apperrors.NewNotFound("Setting")
	}

	return settings[0], nil
}

func (repo *settingRepo) GetByIDAndUser(ctx context.Context, db database.IDB, typ base.SettingType,
	id, userID string, requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	if userID == "" {
		return repo.GetByID(ctx, db, typ, id, requireActive, opts...)
	}
	opts = append(opts, bunex.SelectWhere("setting.object_id = ?", userID))
	return repo.GetByID(ctx, db, typ, id, requireActive, opts...)
}

func (repo *settingRepo) GetByKind(ctx context.Context, db database.IDB, typ base.SettingType, kind string,
	requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	if kind == "" {
		return nil, nil
	}
	setting := &entity.Setting{}
	query := db.NewSelect().Model(setting).
		Where("setting.type = ?", typ).
		Where("setting.kind = ?", kind)
	if requireActive {
		query = query.Where("setting.status = ?", base.SettingStatusActive)
	}
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if setting == nil || errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.NewNotFound("Setting").WithCause(err)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	if hasChange, _ := repo.updateExpiredSetting(ctx, db, setting); hasChange {
		return repo.GetByKind(ctx, db, typ, kind, requireActive, opts...)
	}
	return setting, nil
}

func (repo *settingRepo) GetByName(ctx context.Context, db database.IDB, typ base.SettingType, name string,
	requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	if name == "" {
		return nil, nil
	}
	setting := &entity.Setting{}
	query := db.NewSelect().Model(setting).
		Where("setting.type = ?", typ).
		Where("setting.name = ?", name)
	if requireActive {
		query = query.Where("setting.status = ?", base.SettingStatusActive)
	}
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if setting == nil || errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.NewNotFound("Setting").WithCause(err)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	if hasChange, _ := repo.updateExpiredSetting(ctx, db, setting); hasChange {
		return repo.GetByName(ctx, db, typ, name, requireActive, opts...)
	}
	return setting, nil
}

func (repo *settingRepo) GetByNameAndProject(ctx context.Context, db database.IDB, typ base.SettingType,
	name, projectID string, requireActive bool,
	opts ...bunex.SelectQueryOption) (_ *entity.Setting, err error) {
	if name == "" {
		return nil, nil
	}
	if projectID == "" {
		return repo.GetByName(ctx, db, typ, name, requireActive, opts...)
	}

	opts = append(opts,
		bunex.SelectWhere("setting.name = ?", name),
		bunex.SelectLimit(1),
	)
	if typ != "" {
		opts = append(opts, bunex.SelectWhere("setting.type = ?", typ))
	}
	if requireActive {
		opts = append(opts, bunex.SelectWhere("setting.status = ?", base.SettingStatusActive))
	}

	settings, _, err := repo.ListByProject(ctx, db, projectID, nil, opts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if len(settings) == 0 {
		return nil, apperrors.NewNotFound("Setting")
	}

	return settings[0], nil
}

func (repo *settingRepo) GetByNameAndApp(ctx context.Context, db database.IDB, typ base.SettingType,
	name, projectID, appID string, requireActive bool,
	opts ...bunex.SelectQueryOption) (_ *entity.Setting, err error) {
	if name == "" {
		return nil, nil
	}
	if projectID == "" && appID == "" {
		return repo.GetByName(ctx, db, typ, name, requireActive, opts...)
	}

	opts = append(opts,
		bunex.SelectWhere("setting.name = ?", name),
		bunex.SelectLimit(1),
	)
	if typ != "" {
		opts = append(opts, bunex.SelectWhere("setting.type = ?", typ))
	}
	if requireActive {
		opts = append(opts, bunex.SelectWhere("setting.status = ?", base.SettingStatusActive))
	}

	settings, _, err := repo.ListByApp(ctx, db, projectID, appID, nil, opts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if len(settings) == 0 {
		return nil, apperrors.NewNotFound("Setting")
	}

	return settings[0], nil
}

func (repo *settingRepo) GetByNameAndUser(ctx context.Context, db database.IDB, typ base.SettingType,
	name, userID string, requireActive bool, opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	if userID == "" {
		return repo.GetByName(ctx, db, typ, name, requireActive, opts...)
	}
	opts = append(opts, bunex.SelectWhere("setting.object_id = ?", userID))
	return repo.GetByName(ctx, db, typ, name, requireActive, opts...)
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

func (repo *settingRepo) ListByProject(ctx context.Context, db database.IDB, projectID string,
	paging *basedto.Paging, opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error) {
	var settings []*entity.Setting
	query := db.NewSelect().Model(&settings)

	allOpts := opts
	allOpts = append(allOpts,
		bunex.SelectJoin("LEFT JOIN project_shared_settings pss ON pss.setting_id = setting.id"),
		bunex.SelectWhereGroup(
			bunex.SelectWhere("setting.avail_in_projects = TRUE"),
			bunex.SelectWhereOr("setting.object_id = ?", projectID),
			bunex.SelectWhereOr("(setting.object_id IS NULL AND pss.project_id = ?)", projectID),
		),
	)
	query = bunex.ApplySelect(query, allOpts...)

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
		return repo.ListByProject(ctx, db, projectID, paging, opts...)
	}
	return settings, pagingMeta, nil
}

func (repo *settingRepo) ListByApp(ctx context.Context, db database.IDB, projectID, appID string,
	paging *basedto.Paging, opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error) {
	var settings []*entity.Setting
	query := db.NewSelect().Model(&settings)

	allOpts := opts
	if projectID != "" {
		allOpts = append(allOpts,
			bunex.SelectJoin("LEFT JOIN project_shared_settings pss ON pss.setting_id = setting.id"),
			bunex.SelectWhereGroup(
				bunex.SelectWhere("setting.object_id = ?", appID),
				bunex.SelectWhereOrGroup(
					bunex.SelectWhereOr("setting.avail_in_projects = TRUE"),
					bunex.SelectWhereOr("setting.object_id = ?", projectID),
					bunex.SelectWhereOr("(setting.object_id IS NULL AND pss.project_id = ?)", projectID),
				),
			),
		)
	} else {
		allOpts = append(allOpts,
			bunex.SelectWhere("setting.object_id = ?", appID),
		)
	}
	query = bunex.ApplySelect(query, allOpts...)

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
		return repo.ListByApp(ctx, db, projectID, appID, paging, opts...)
	}
	return settings, pagingMeta, nil
}

func (repo *settingRepo) ListByUser(ctx context.Context, db database.IDB, userID string,
	paging *basedto.Paging, opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error) {
	if userID == "" {
		return repo.List(ctx, db, paging, opts...)
	}
	opts = append(opts, bunex.SelectWhere("setting.object_id = ?", userID))
	return repo.List(ctx, db, paging, opts...)
}

func (repo *settingRepo) ListByIDs(ctx context.Context, db database.IDB, ids []string, requireActive bool,
	opts ...bunex.SelectQueryOption) ([]*entity.Setting, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var settings []*entity.Setting
	query := db.NewSelect().Model(&settings).Where("setting.id IN (?)", bun.In(ids))
	if requireActive {
		query = query.Where("setting.status = ?", base.SettingStatusActive)
	}
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if err != nil {
		return nil, apperrors.New(err)
	}

	if hasChange, _ := repo.updateExpiredSettings(ctx, db, settings); hasChange {
		return repo.ListByIDs(ctx, db, ids, requireActive, opts...)
	}
	return settings, nil
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

func (repo *settingRepo) updateExpiredSetting(ctx context.Context, db database.IDB, setting *entity.Setting) (
	bool, error) {
	if setting == nil {
		return false, nil
	}
	return repo.updateExpiredSettings(ctx, db, []*entity.Setting{setting})
}

func (repo *settingRepo) updateExpiredSettings(ctx context.Context, db database.IDB, settings []*entity.Setting) (
	hasChange bool, err error) {
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
