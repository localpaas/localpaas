package containeragentuc

import (
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/services/docker"
)

type UC struct {
	logger        logging.Logger
	db            *database.DB
	dockerManager docker.Manager
}

func New(
	logger logging.Logger,
	db *database.DB,
	dockerManager docker.Manager,
) *UC {
	return &UC{
		logger:        logger,
		db:            db,
		dockerManager: dockerManager,
	}
}
