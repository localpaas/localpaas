package oauthdto

import (
	"github.com/localpaas/localpaas/localpaas_app/base"
)

type GetOAuthNoAuthReq struct {
	ID     string
	Kind   string
	Status []base.SettingStatus
}
