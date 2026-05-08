package sysbackupservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type Service interface {
	Backup(ctx context.Context, db database.Tx, req *SysBackupReq) (*SysBackupResp, error)
}
