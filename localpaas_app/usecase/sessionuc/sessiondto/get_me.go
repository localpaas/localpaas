package sessiondto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/useruc/userdto"
)

type GetMeReq struct {
	GetAccesses bool `json:"-" mapstructure:"getAccesses"`
}

func NewGetMeReq() *GetMeReq {
	return &GetMeReq{}
}

type GetMeResp struct {
	Meta *basedto.Meta  `json:"meta"`
	Data *GetMeDataResp `json:"data"`
}

type GetMeDataResp struct {
	NextStep string                   `json:"nextStep,omitempty"`
	User     *userdto.UserDetailsResp `json:"user"`
}

func TransformUserDetails(user *entity.User) (resp *userdto.UserDetailsResp, err error) {
	resp, err = userdto.TransformUserDetails(user)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
