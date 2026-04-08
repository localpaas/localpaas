package appservice

import (
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type PersistingAppData struct {
	UpsertingApps        []*entity.App
	UpsertingTags        []*entity.AppTag
	UpsertingSettings    []*entity.Setting
	UpsertingDeployments []*entity.Deployment
	UpsertingTasks       []*entity.Task

	AppsToDeleteTags []string
}
