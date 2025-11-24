package lpappuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/service/lpappservice"
)

type LpAppUC struct {
	db           *database.DB
	lpAppService lpappservice.LpAppService
}

func NewLpAppUC(
	db *database.DB,
	lpAppService lpappservice.LpAppService,
) *LpAppUC {
	return &LpAppUC{
		db:           db,
		lpAppService: lpAppService,
	}
}
