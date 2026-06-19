package imagebuildservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type Service interface {
	ImageBuild(ctx context.Context, db database.IDB, req *ImageBuildReq) (*ImageBuildResp, error)
}
