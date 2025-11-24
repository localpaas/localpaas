package lpappdto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type RestartLpAppReq struct {
}

func NewRestartLpAppReq() *RestartLpAppReq {
	return &RestartLpAppReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *RestartLpAppReq) Validate() apperrors.ValidationErrors {
	return nil
}

type RestartLpAppResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
