package nginxuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/service/nginxservice"
)

type NginxUC struct {
	db           *database.DB
	nginxService nginxservice.NginxService
}

func NewNginxUC(
	db *database.DB,
	nginxService nginxservice.NginxService,
) *NginxUC {
	return &NginxUC{
		db:           db,
		nginxService: nginxService,
	}
}
