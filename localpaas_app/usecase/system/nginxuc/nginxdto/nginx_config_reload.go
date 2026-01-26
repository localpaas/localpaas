package nginxdto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type ReloadNginxConfigReq struct {
}

func NewReloadNginxConfigReq() *ReloadNginxConfigReq {
	return &ReloadNginxConfigReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *ReloadNginxConfigReq) Validate() apperrors.ValidationErrors {
	return nil
}

type ReloadNginxConfigResp struct {
	Meta *basedto.Meta `json:"meta"`
}
