package lpappdto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type RestartLpAppReq struct {
	RestartMainApp  bool `json:"restartMainApp"`
	RestartDbApp    bool `json:"restartDbApp"`
	RestartCacheApp bool `json:"restartCacheApp"`
}

func NewRestartLpAppReq() *RestartLpAppReq {
	return &RestartLpAppReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *RestartLpAppReq) Validate() apperrors.ValidationErrors {
	return nil
}

type RestartLpAppResp struct {
	Meta *basedto.Meta `json:"meta"`
}
