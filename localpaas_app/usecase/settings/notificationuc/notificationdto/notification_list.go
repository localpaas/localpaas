package notificationdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ListNotificationReq struct {
	settings.ListSettingReq
}

func NewListNotificationReq() *ListNotificationReq {
	return &ListNotificationReq{
		ListSettingReq: settings.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListNotificationReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListNotificationResp struct {
	Meta *basedto.ListMeta   `json:"meta"`
	Data []*NotificationResp `json:"data"`
}

func TransformNotifications(
	settings []*entity.Setting,
	refObjects *entity.RefObjects,
) (resp []*NotificationResp, err error) {
	resp = make([]*NotificationResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformNotification(setting, refObjects)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
