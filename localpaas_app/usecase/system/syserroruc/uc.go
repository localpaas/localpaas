package syserroruc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
)

type SysErrorUC struct {
	db           *database.DB
	appErrorRepo repository.SysErrorRepo
}

func NewSysErrorUC(
	db *database.DB,
	appErrorRepo repository.SysErrorRepo,
) *SysErrorUC {
	return &SysErrorUC{
		db:           db,
		appErrorRepo: appErrorRepo,
	}
}
