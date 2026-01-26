package nginxdto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type ResetNginxConfigReq struct {
}

func NewResetNginxConfigReq() *ResetNginxConfigReq {
	return &ResetNginxConfigReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *ResetNginxConfigReq) Validate() apperrors.ValidationErrors {
	return nil
}

type ResetNginxConfigResp struct {
	Meta *basedto.Meta `json:"meta"`
}
