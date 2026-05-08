package sysbackupservice

import (
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/tasks/queue"
)

type SysBackupReq struct {
	*queue.TaskExecData
	CronJob *entity.Setting
}

type SysBackupResp struct {
	SkipResultNotification bool
}
