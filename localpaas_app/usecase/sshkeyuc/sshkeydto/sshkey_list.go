package sshkeydto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type ListSSHKeyReq struct {
	Search string `json:"-" mapstructure:"search"`

	Paging basedto.Paging `json:"-"`
}

func NewListSSHKeyReq() *ListSSHKeyReq {
	return &ListSSHKeyReq{
		Paging: basedto.Paging{
			// Default paging if unset by client
			Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
		},
	}
}

func (req *ListSSHKeyReq) Validate() apperrors.ValidationErrors {
	return nil
}

type ListSSHKeyResp struct {
	Meta *basedto.Meta `json:"meta"`
	Data []*SSHKeyResp `json:"data"`
}

func TransformSSHKeys(settings []*entity.Setting) ([]*SSHKeyResp, error) {
	return basedto.TransformObjectSlice(settings, TransformSSHKey) //nolint:wrapcheck
}
