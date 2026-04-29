package lpappdto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/service/lpappservice"
)

type GetLpAppReleaseInfoReq struct {
}

func NewGetLpAppReleaseInfoReq() *GetLpAppReleaseInfoReq {
	return &GetLpAppReleaseInfoReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *GetLpAppReleaseInfoReq) Validate() apperrors.ValidationErrors {
	return nil
}

type GetLpAppReleaseInfoResp struct {
	Meta *basedto.Meta         `json:"meta"`
	Data *LpAppReleaseInfoResp `json:"data"`
}

type LpAppReleaseInfoResp struct {
	*lpappservice.AppReleaseInfo
}
