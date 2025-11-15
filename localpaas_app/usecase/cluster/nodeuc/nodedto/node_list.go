package nodedto

import (
	"github.com/docker/docker/api/types/swarm"
	vld "github.com/tiendc/go-validator"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type ListNodeReq struct {
	Status []base.NodeStatus `json:"-" mapstructure:"status"`
	Role   []base.NodeRole   `json:"-" mapstructure:"role"`
	Search string            `json:"-" mapstructure:"search"`

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
	validators = append(validators, basedto.ValidateSlice(req.Role, true, 0,
		base.AllNodeRoles, "role")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListNodeResp struct {
	Meta *basedto.Meta `json:"meta"`
	Data []*NodeResp   `json:"data"`
}

func TransformNodes(nodes []swarm.Node, detailed bool) []*NodeResp {
	return gofn.MapSlice(nodes, func(node swarm.Node) *NodeResp {
		return TransformNode(&node, detailed)
	})
}
