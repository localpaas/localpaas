package oauthdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type ListOAuthReq struct {
	settings.ListSettingReq
	Kind []base.OAuthKind `json:"-" mapstructure:"kind"`
}

func NewListOAuthReq() *ListOAuthReq {
	return &ListOAuthReq{
		ListSettingReq: settings.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListOAuthReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	validators = append(validators, basedto.ValidateSlice(req.Kind, true, 0, base.AllOAuthKinds, "kind")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListOAuthResp struct {
	Meta *basedto.ListMeta `json:"meta"`
	Data []*OAuthResp      `json:"data"`
}

func TransformOAuths(settings []*entity.Setting, baseCallbackURL string, objectID string) ([]*OAuthResp, error) {
	resp := make([]*OAuthResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformOAuth(setting, baseCallbackURL, objectID)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
