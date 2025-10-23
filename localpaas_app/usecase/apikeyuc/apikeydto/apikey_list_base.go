package apikeydto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type ListAPIKeyBaseReq struct {
	Search string `json:"-" mapstructure:"search"`

	Paging basedto.Paging `json:"-"`
}

func NewListAPIKeyBaseReq() *ListAPIKeyBaseReq {
	return &ListAPIKeyBaseReq{
		Paging: basedto.Paging{
			// Default paging if unset by client
			Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
		},
	}
}

// Validate implements interface basedto.ReqValidator
func (req *ListAPIKeyBaseReq) Validate() apperrors.ValidationErrors {
	return nil
}

type ListAPIKeyBaseResp struct {
	Meta *basedto.Meta     `json:"meta"`
	Data []*APIKeyBaseResp `json:"data"`
}
