package clusterdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type DeleteNodeReq struct {
	NodeID string `json:"-"`
}

func NewDeleteNodeReq() *DeleteNodeReq {
	return &DeleteNodeReq{}
}

func (req *DeleteNodeReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.NodeID, true, "nodeId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteNodeResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
