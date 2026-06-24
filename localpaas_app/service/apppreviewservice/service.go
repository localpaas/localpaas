package apppreviewservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type Service interface {
	CreatePreview(ctx context.Context, db database.Tx, req *CreatePreviewReq) (*CreatePreviewResp, error)
}
