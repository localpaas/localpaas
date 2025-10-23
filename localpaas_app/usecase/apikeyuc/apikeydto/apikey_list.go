package apikeydto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type ListAPIKeyReq struct {
	Search string `json:"-" mapstructure:"search"`

	Paging basedto.Paging `json:"-"`
}

func NewListAPIKeyReq() *ListAPIKeyReq {
	return &ListAPIKeyReq{
		Paging: basedto.Paging{
			// Default paging if unset by client
			Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
		},
	}
}

func (req *ListAPIKeyReq) Validate() apperrors.ValidationErrors {
	return nil
}

type ListAPIKeyResp struct {
	Meta *basedto.Meta `json:"meta"`
	Data []*APIKeyResp `json:"data"`
}

func TransformAPIKeys(settings []*entity.Setting, userMap map[string]*entity.User) ([]*APIKeyResp, error) {
	//nolint:wrapcheck
	return basedto.TransformObjectSlice(settings, func(setting *entity.Setting) (*APIKeyResp, error) {
		return TransformAPIKey(setting, userMap)
	})
}
