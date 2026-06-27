package apppreviewservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
)

type Service interface {
	GetPreview(ctx context.Context, db database.IDB, appID, repoRef string, extraOpts ...bunex.SelectQueryOption) (
		*entity.App, error)
	GetPreviews(ctx context.Context, db database.IDB, appID string, extraOpts ...bunex.SelectQueryOption) (
		[]*entity.App, error)

	CreatePreview(ctx context.Context, db database.Tx, req *CreatePreviewReq) (*CreatePreviewResp, error)
}
