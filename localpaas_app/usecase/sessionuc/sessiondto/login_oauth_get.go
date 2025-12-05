package sessiondto

import (
	"github.com/localpaas/localpaas/localpaas_app/base"
)

type GetLoginOAuthReq struct {
	ID     string
	Kind   string
	Status []base.SettingStatus
}
