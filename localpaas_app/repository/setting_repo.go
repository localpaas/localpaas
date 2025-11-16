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
	GetByID(ctx context.Context, db database.IDB, id string,
		opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetByKind(ctx context.Context, db database.IDB, typ base.SettingType, kind string,
		opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	GetByName(ctx context.Context, db database.IDB, typ base.SettingType, name string,
		opts ...bunex.SelectQueryOption) (*entity.Setting, error)
	List(ctx context.Context, db database.IDB, paging *basedto.Paging,
		opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error)
	ListByIDs(ctx context.Context, db database.IDB, ids []string,
		opts ...bunex.SelectQueryOption) ([]*entity.Setting, error)

	Upsert(ctx context.Context, db database.IDB, setting *entity.Setting,
		conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error
	UpsertMulti(ctx context.Context, db database.IDB, settings []*entity.Setting,
		conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error
	Update(ctx context.Context, db database.IDB, setting *entity.Setting,
		opts ...bunex.UpdateQueryOption) error
	UpdateMulti(ctx context.Context, db database.IDB, settings []*entity.Setting,
		opts ...bunex.UpdateQueryOption) error
}

type settingRepo struct {
}

func NewSettingRepo() SettingRepo {
	return &settingRepo{}
}

func (repo *settingRepo) GetByID(ctx context.Context, db database.IDB, id string,
	opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	setting := &entity.Setting{}
	query := db.NewSelect().Model(setting).Where("setting.id = ?", id)
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if setting == nil || errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.NewNotFound("Setting").WithCause(err)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	if hasChange, _ := repo.updateExpiredSetting(ctx, db, setting); hasChange {
		return repo.GetByID(ctx, db, id, opts...)
	}
	return setting, nil
}

func (repo *settingRepo) GetByKind(ctx context.Context, db database.IDB, typ base.SettingType, kind string,
	opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	if kind == "" {
		return nil, nil
	}
	setting := &entity.Setting{}
	query := db.NewSelect().Model(setting).
		Where("setting.type = ?", typ).
		Where("setting.kind = ?", kind)
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if setting == nil || errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.NewNotFound("Setting").WithCause(err)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	if hasChange, _ := repo.updateExpiredSetting(ctx, db, setting); hasChange {
		return repo.GetByKind(ctx, db, typ, kind, opts...)
	}
	return setting, nil
}

func (repo *settingRepo) GetByName(ctx context.Context, db database.IDB, typ base.SettingType, name string,
	opts ...bunex.SelectQueryOption) (*entity.Setting, error) {
	if name == "" {
		return nil, nil
	}
	setting := &entity.Setting{}
	query := db.NewSelect().Model(setting).
		Where("setting.type = ?", typ).
		Where("setting.name = ?", name)
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if setting == nil || errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.NewNotFound("Setting").WithCause(err)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	if hasChange, _ := repo.updateExpiredSetting(ctx, db, setting); hasChange {
		return repo.GetByName(ctx, db, typ, name, opts...)
	}
	return setting, nil
}

func (repo *settingRepo) List(ctx context.Context, db database.IDB, paging *basedto.Paging,
	opts ...bunex.SelectQueryOption) ([]*entity.Setting, *basedto.PagingMeta, error) {
	var settings []*entity.Setting
	query := db.NewSelect().Model(&settings)
	query = bunex.ApplySelect(query, opts...)

	pagingMeta := newPagingMeta(paging)

	// Counts the total first
	if paging != nil {
		total, err := query.Count(ctx)
		if err != nil {
			return nil, nil, apperrors.Wrap(err)
		}
		pagingMeta.Total = total
	}

	// Applies pagination
	query = bunex.ApplyPagination(query, paging)
	err := query.Scan(ctx)
	if err != nil {
		return nil, nil, wrapPaginationError(err, paging)
	}

	if hasChange, _ := repo.updateExpiredSettings(ctx, db, settings); hasChange {
		return repo.List(ctx, db, paging, opts...)
	}
	return settings, pagingMeta, nil
}

func (repo *settingRepo) ListByIDs(ctx context.Context, db database.IDB, ids []string,
	opts ...bunex.SelectQueryOption) ([]*entity.Setting, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var settings []*entity.Setting
	query := db.NewSelect().Model(&settings).Where("setting.id IN (?)", bun.In(ids))
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if err != nil {
		return nil, apperrors.New(err)
	}

	if hasChange, _ := repo.updateExpiredSettings(ctx, db, settings); hasChange {
		return repo.ListByIDs(ctx, db, ids, opts...)
	}
	return settings, nil
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
	return repo.UpdateMulti(ctx, db, []*entity.Setting{setting}, opts...)
}

func (repo *settingRepo) UpdateMulti(ctx context.Context, db database.IDB, settings []*entity.Setting,
	opts ...bunex.UpdateQueryOption) error {
	if len(settings) == 0 {
		return nil
	}

	query := db.NewUpdate().Model(&settings).WherePK()
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
