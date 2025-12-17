package nodedto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type JoinNodeReq struct {
	Host          string              `json:"host"`
	Port          int                 `json:"port"`
	User          string              `json:"user"`
	SSHKey        basedto.ObjectIDReq `json:"sshKey"`
	JoinAsManager bool                `json:"joinAsManager"`
}

func NewJoinNodeReq() *JoinNodeReq {
	return &JoinNodeReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *JoinNodeReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateObjectIDReq(&req.SSHKey, true, "sshKey")...)
	validators = append(validators, basedto.ValidateStr(&req.Host, true, 1, 100, "host")...)      //nolint:mnd
	validators = append(validators, basedto.ValidateStr(&req.User, true, 1, 100, "user")...)      //nolint:mnd
	validators = append(validators, basedto.ValidateNumber(&req.Port, true, 1, 65535, "port")...) //nolint:mnd
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type JoinNodeResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *JoinNodeDataResp `json:"data"`
}

type JoinNodeDataResp struct {
	CommandOutput string `json:"commandOutput"`
}
