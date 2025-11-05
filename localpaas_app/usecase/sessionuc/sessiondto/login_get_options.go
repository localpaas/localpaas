package sessiondto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
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
	Meta *basedto.BaseMeta `json:"meta"`
	Data *LoginOptionsResp `json:"data"`
}

type LoginOptionsResp struct {
	AllowLoginWithGitHub bool `json:"allowLoginWithGitHub"`
	AllowLoginWithGitLab bool `json:"allowLoginWithGitLab"`
}
