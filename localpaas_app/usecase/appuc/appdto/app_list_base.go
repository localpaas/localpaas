package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type ListAppBaseReq struct {
	ProjectID    string           `json:"-" mapstructure:"projectId"`
	ProjectEnvID string           `json:"-" mapstructure:"projectEnvId"`
	Status       []base.AppStatus `json:"-" mapstructure:"status"`
	Search       string           `json:"-" mapstructure:"search"`

	Paging basedto.Paging `json:"-"`
}

func NewListAppBaseReq() *ListAppBaseReq {
	return &ListAppBaseReq{
		Status: []base.AppStatus{base.AppStatusActive},
		Paging: basedto.Paging{
			// Default paging if unset by client
			Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
		},
	}
}

// Validate implements interface basedto.ReqValidator
func (req *ListAppBaseReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, false, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.ProjectEnvID, false, "projectEnvId")...)
	validators = append(validators, basedto.ValidateSlice(req.Status, true, 0,
		base.AllAppStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListAppBaseResp struct {
	Meta *basedto.Meta  `json:"meta"`
	Data []*AppBaseResp `json:"data"`
}
