package cacheentity

import "github.com/localpaas/localpaas/localpaas_app/base"

type TaskControl struct {
	ID  string           `json:"id"`
	Cmd base.TaskCommand `json:"cmd"`
}
