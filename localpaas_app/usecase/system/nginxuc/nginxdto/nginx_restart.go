package nginxdto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type RestartNginxReq struct {
}

func NewRestartNginxReq() *RestartNginxReq {
	return &RestartNginxReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *RestartNginxReq) Validate() apperrors.ValidationErrors {
	return nil
}

type RestartNginxResp struct {
	Meta *basedto.Meta `json:"meta"`
}
