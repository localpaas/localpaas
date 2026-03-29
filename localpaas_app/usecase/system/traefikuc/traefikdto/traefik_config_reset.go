package traefikdto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type ResetTraefikConfigReq struct {
}

func NewResetTraefikConfigReq() *ResetTraefikConfigReq {
	return &ResetTraefikConfigReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *ResetTraefikConfigReq) Validate() apperrors.ValidationErrors {
	return nil
}

type ResetTraefikConfigResp struct {
	Meta *basedto.Meta `json:"meta"`
}
