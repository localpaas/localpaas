package sshkeydto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ListSSHKeyReq struct {
	settings.ListSettingReq
}

func NewListSSHKeyReq() *ListSSHKeyReq {
	return &ListSSHKeyReq{
		ListSettingReq: settings.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListSSHKeyReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListSSHKeyResp struct {
	Meta *basedto.Meta `json:"meta"`
	Data []*SSHKeyResp `json:"data"`
}

func TransformSSHKeys(settings []*entity.Setting, objectID string) (resp []*SSHKeyResp, err error) {
	resp = make([]*SSHKeyResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformSSHKey(setting, objectID)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
