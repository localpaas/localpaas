package traefikdto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type ReloadTraefikConfigReq struct {
}

func NewReloadTraefikConfigReq() *ReloadTraefikConfigReq {
	return &ReloadTraefikConfigReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *ReloadTraefikConfigReq) Validate() apperrors.ValidationErrors {
	return nil
}

type ReloadTraefikConfigResp struct {
	Meta *basedto.Meta `json:"meta"`
}
