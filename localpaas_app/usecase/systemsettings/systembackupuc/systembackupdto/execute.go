package systembackupdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ExecuteSystemBackupReq struct {
	settings.GetSettingReq
}

func NewExecuteSystemBackupReq() *ExecuteSystemBackupReq {
	return &ExecuteSystemBackupReq{}
}

func (req *ExecuteSystemBackupReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ExecuteSystemBackupResp struct {
	Meta *basedto.Meta                `json:"meta"`
	Data *ExecuteSystemBackupDataResp `json:"data"`
}

type ExecuteSystemBackupDataResp struct {
	Task *basedto.ObjectIDResp `json:"task"`
}
