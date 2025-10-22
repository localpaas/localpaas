package projectenvdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type ListProjectEnvBaseReq struct {
	ProjectID string               `json:"-"`
	Status    []base.ProjectStatus `json:"-" mapstructure:"status"`
	Search    string               `json:"-" mapstructure:"search"`

	Paging basedto.Paging `json:"-"`
}

func NewListProjectEnvBaseReq() *ListProjectEnvBaseReq {
	return &ListProjectEnvBaseReq{
		Status: []base.ProjectStatus{base.ProjectStatusActive},
		Paging: basedto.Paging{
			// Default paging if unset by client
			Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
		},
	}
}

// Validate implements interface basedto.ReqValidator
func (req *ListProjectEnvBaseReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateSlice(req.Status, true, 0,
		base.AllProjectStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListProjectEnvBaseResp struct {
	Meta *basedto.Meta         `json:"meta"`
	Data []*ProjectEnvBaseResp `json:"data"`
}
