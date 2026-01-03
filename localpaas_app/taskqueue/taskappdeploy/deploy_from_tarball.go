package taskappdeploy

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

func (e *Executor) deployFromTarball(
	ctx context.Context,
	db database.Tx,
	taskData *taskData,
) error {
	// TODO
	return nil
}
