package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
)

type NodeRepo interface {
	GetByID(ctx context.Context, db database.IDB, id string,
		opts ...bunex.SelectQueryOption) (*entity.Node, error)
	GetByHostName(ctx context.Context, db database.IDB, hostname string,
		opts ...bunex.SelectQueryOption) (*entity.Node, error)
	GetByIP(ctx context.Context, db database.IDB, ip string,
		opts ...bunex.SelectQueryOption) (*entity.Node, error)
	List(ctx context.Context, db database.IDB, paging *basedto.Paging,
		opts ...bunex.SelectQueryOption) ([]*entity.Node, *basedto.PagingMeta, error)
	ListByIDs(ctx context.Context, db database.IDB, ids []string,
		opts ...bunex.SelectQueryOption) ([]*entity.Node, error)

	Upsert(ctx context.Context, db database.IDB, node *entity.Node,
		conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error
	UpsertMulti(ctx context.Context, db database.IDB, nodes []*entity.Node,
		conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error
}

type nodeRepo struct {
}

func NewNodeRepo() NodeRepo {
	return &nodeRepo{}
}

func (repo *nodeRepo) GetByID(ctx context.Context, db database.IDB, id string,
	opts ...bunex.SelectQueryOption) (*entity.Node, error) {
	node := &entity.Node{}
	query := db.NewSelect().Model(node).Where("node.id = ?", id)
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if node == nil || errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.NewNotFound("Node").WithCause(err)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return node, nil
}

func (repo *nodeRepo) GetByHostName(ctx context.Context, db database.IDB, hostname string,
	opts ...bunex.SelectQueryOption) (*entity.Node, error) {
	node := &entity.Node{}
	query := db.NewSelect().Model(node).Where("node.host_name = ?", hostname)
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if node == nil || errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.NewNotFound("Node").WithCause(err)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return node, nil
}

func (repo *nodeRepo) GetByIP(ctx context.Context, db database.IDB, ip string,
	opts ...bunex.SelectQueryOption) (*entity.Node, error) {
	node := &entity.Node{}
	query := db.NewSelect().Model(node).Where("node.ip = ?", ip)
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if node == nil || errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.NewNotFound("Node").WithCause(err)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return node, nil
}

func (repo *nodeRepo) List(ctx context.Context, db database.IDB, paging *basedto.Paging,
	opts ...bunex.SelectQueryOption) ([]*entity.Node, *basedto.PagingMeta, error) {
	var nodes []*entity.Node
	query := db.NewSelect().Model(&nodes)
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

	return nodes, pagingMeta, nil
}

func (repo *nodeRepo) ListByIDs(ctx context.Context, db database.IDB, ids []string,
	opts ...bunex.SelectQueryOption) ([]*entity.Node, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var nodes []*entity.Node
	query := db.NewSelect().Model(&nodes).Where("node.id IN (?)", bun.In(ids))
	query = bunex.ApplySelect(query, opts...)

	err := query.Scan(ctx)
	if err != nil {
		return nil, apperrors.New(err)
	}
	return nodes, nil
}

func (repo *nodeRepo) Upsert(ctx context.Context, db database.IDB, node *entity.Node,
	conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error {
	return repo.UpsertMulti(ctx, db, []*entity.Node{node}, conflictCols, updateCols, opts...)
}

func (repo *nodeRepo) UpsertMulti(ctx context.Context, db database.IDB, nodes []*entity.Node,
	conflictCols, updateCols []string, opts ...bunex.InsertQueryOption) error {
	if len(nodes) == 0 {
		return nil
	}
	query := db.NewInsert().Model(&nodes)
	query = bunex.ApplyInsert(query, opts...)
	query = bunex.ApplyUpsert(query, conflictCols, updateCols)

	_, err := query.Exec(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
