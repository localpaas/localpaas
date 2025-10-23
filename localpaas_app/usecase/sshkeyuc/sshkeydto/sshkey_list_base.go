package sshkeydto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type ListSSHKeyBaseReq struct {
	Search string `json:"-" mapstructure:"search"`

	Paging basedto.Paging `json:"-"`
}

func NewListSSHKeyBaseReq() *ListSSHKeyBaseReq {
	return &ListSSHKeyBaseReq{
		Paging: basedto.Paging{
			// Default paging if unset by client
			Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
		},
	}
}

// Validate implements interface basedto.ReqValidator
func (req *ListSSHKeyBaseReq) Validate() apperrors.ValidationErrors {
	return nil
}

type ListSSHKeyBaseResp struct {
	Meta *basedto.Meta     `json:"meta"`
	Data []*SSHKeyBaseResp `json:"data"`
}
