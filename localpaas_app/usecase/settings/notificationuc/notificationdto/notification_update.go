package notificationdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateNotificationReq struct {
	settings.UpdateSettingReq
	*NotificationBaseReq
}

func NewUpdateNotificationReq() *UpdateNotificationReq {
	return &UpdateNotificationReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateNotificationReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateNotificationResp struct {
	Meta *basedto.Meta `json:"meta"`
}
