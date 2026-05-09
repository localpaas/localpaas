package healthcheckservice

import (
	"github.com/localpaas/localpaas/localpaas_app/tasks/queue"
)

type HealthcheckReq struct {
	*queue.HealthcheckExecData
}

type HealthcheckResp struct {
}
