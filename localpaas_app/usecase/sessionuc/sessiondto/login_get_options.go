package sessiondto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type GetLoginOptionsReq struct {
}

func NewGetLoginOptionsReq() *GetLoginOptionsReq {
	return &GetLoginOptionsReq{}
}

func (req *GetLoginOptionsReq) Validate() apperrors.ValidationErrors {
	return nil
}

type GetLoginOptionsResp struct {
	Meta *basedto.BaseMeta  `json:"meta"`
	Data []*LoginOptionResp `json:"data"`
}

type LoginOptionResp struct {
	Type    base.OAuthType `json:"type"`
	Name    string         `json:"name"`
	Icon    string         `json:"icon"`
	AuthURL string         `json:"authURL"`
}
