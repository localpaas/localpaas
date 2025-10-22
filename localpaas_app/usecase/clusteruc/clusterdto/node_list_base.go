package clusterdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type ListNodeBaseReq struct {
	Status      []base.NodeStatus `json:"-" mapstructure:"status"`
	InfraStatus string            `json:"-" mapstructure:"infraStatus"`
	Search      string            `json:"-" mapstructure:"search"`

	Paging basedto.Paging `json:"-"`
}

func NewListNodeBaseReq() *ListNodeBaseReq {
	return &ListNodeBaseReq{
		Status: []base.NodeStatus{base.NodeStatusActive},
		Paging: basedto.Paging{
			// Default paging if unset by client
			Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "host_name"}},
		},
	}
}

// Validate implements interface basedto.ReqValidator
func (req *ListNodeBaseReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateSlice(req.Status, true, 0,
		base.AllNodeStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListNodeBaseResp struct {
	Meta *basedto.Meta   `json:"meta"`
	Data []*NodeBaseResp `json:"data"`
}
