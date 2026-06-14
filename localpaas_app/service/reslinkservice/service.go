package reslinkservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type Service interface {
	SetLinks(ctx context.Context, db database.IDB, srcType base.ResourceType, srcID string,
		dstType base.ResourceType, dstIDs []string, options ...SetLinkOption) error
	AddLinks(ctx context.Context, db database.IDB, srcType base.ResourceType, srcID string,
		dstType base.ResourceType, dstIDs []string, options ...SetLinkOption) error
	RemoveLinks(ctx context.Context, db database.IDB, srcType base.ResourceType, srcID string,
		dstType base.ResourceType, dstIDs []string) error
}
