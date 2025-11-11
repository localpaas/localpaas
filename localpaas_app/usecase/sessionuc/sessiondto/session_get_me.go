package sessiondto

import (
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type GetMeReq struct {
}

func NewGetMeReq() *GetMeReq {
	return &GetMeReq{}
}

type GetMeResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *GetMeDataResp    `json:"data"`
}

type GetMeDataResp struct {
	NextStep string    `json:"nextStep,omitempty"`
	User     *UserResp `json:"user"`
}
