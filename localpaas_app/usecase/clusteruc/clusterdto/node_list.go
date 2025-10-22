package clusterdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type ListNodeReq struct {
	Status      []base.NodeStatus `json:"-" mapstructure:"status"`
	InfraStatus []string          `json:"-" mapstructure:"infraStatus"`
	Search      string            `json:"-" mapstructure:"search"`

	Paging basedto.Paging `json:"-"`
}

func NewListNodeReq() *ListNodeReq {
	return &ListNodeReq{
		Paging: basedto.Paging{
			// Default paging if unset by client
			Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "host_name"}},
		},
	}
}

func (req *ListNodeReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateSlice(req.Status, true, 0,
		base.AllNodeStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListNodeResp struct {
	Meta *basedto.Meta `json:"meta"`
	Data []*NodeResp   `json:"data"`
}

func TransformNodes(nodes []*entity.Node) ([]*NodeResp, error) {
	return basedto.TransformObjectSlice(nodes, TransformNode) //nolint:wrapcheck
}
