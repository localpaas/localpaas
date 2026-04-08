package clusterservice

import "github.com/localpaas/localpaas/localpaas_app/entity"

type PersistingClusterData struct {
	UpsertingSettings []*entity.Setting
}
