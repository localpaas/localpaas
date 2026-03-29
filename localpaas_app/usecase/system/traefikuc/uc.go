package traefikuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/service/traefikservice"
)

type TraefikUC struct {
	db             *database.DB
	traefikService traefikservice.TraefikService
}

func NewTraefikUC(
	db *database.DB,
	traefikService traefikservice.TraefikService,
) *TraefikUC {
	return &TraefikUC{
		db:             db,
		traefikService: traefikService,
	}
}
