package lpappdto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type ReloadLpAppConfigReq struct {
}

func NewReloadLpAppConfigReq() *ReloadLpAppConfigReq {
	return &ReloadLpAppConfigReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *ReloadLpAppConfigReq) Validate() apperrors.ValidationErrors {
	return nil
}

type ReloadLpAppConfigResp struct {
	Meta *basedto.Meta `json:"meta"`
}
