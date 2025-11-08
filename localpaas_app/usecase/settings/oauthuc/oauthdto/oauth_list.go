package oauthdto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type ListOAuthReq struct {
	Search string `json:"-" mapstructure:"search"`

	Paging basedto.Paging `json:"-"`
}

func NewListOAuthReq() *ListOAuthReq {
	return &ListOAuthReq{
		Paging: basedto.Paging{
			// Default paging if unset by client
			Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
		},
	}
}

func (req *ListOAuthReq) Validate() apperrors.ValidationErrors {
	return nil
}

type ListOAuthResp struct {
	Meta *basedto.Meta `json:"meta"`
	Data []*OAuthResp  `json:"data"`
}

func TransformOAuths(settings []*entity.Setting) ([]*OAuthResp, error) {
	resp, err := basedto.TransformObjectSlice(settings, TransformOAuth)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
