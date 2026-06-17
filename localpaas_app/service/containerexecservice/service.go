package containerexecservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type Service interface {
	ContainerExec(ctx context.Context, req *ContainerExecReq) (*ContainerExecResp, error)
	SchedJobExec(ctx context.Context, db database.Tx, req *SchedJobExecReq) (*SchedJobExecResp, error)
}
