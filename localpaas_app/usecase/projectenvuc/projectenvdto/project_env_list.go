package projectenvdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type ListProjectEnvReq struct {
	ProjectID string               `json:"-"`
	Status    []base.ProjectStatus `json:"-" mapstructure:"status"`
	Search    string               `json:"-" mapstructure:"search"`

	Paging basedto.Paging `json:"-"`
}

func NewListProjectEnvReq() *ListProjectEnvReq {
	return &ListProjectEnvReq{
		Paging: basedto.Paging{
			// Default paging if unset by client
			Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
		},
	}
}

func (req *ListProjectEnvReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateSlice(req.Status, true, 0,
		base.AllProjectStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListProjectEnvResp struct {
	Meta *basedto.Meta     `json:"meta"`
	Data []*ProjectEnvResp `json:"data"`
}

func TransformProjectEnvs(projectEnvs []*entity.ProjectEnv) ([]*ProjectEnvResp, error) {
	return basedto.TransformObjectSlice(projectEnvs, TransformProjectEnv) //nolint:wrapcheck
}
