package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
)

//
// REQUEST
//

type HttpSettingsReq struct {
	Enabled bool `json:"enabled"`
}

// nolint
func (req *HttpSettingsReq) validate(field string) (res []vld.Validator) {
	if req == nil {
		return
	}
	if field != "" {
		field += "."
	}
	// TODO:
	return res
}

//
// RESPONSE
//

type HttpSettingsResp struct {
	Enabled bool `json:"enabled"`
}

func TransformHttpSettings(setting *entity.Setting) (resp *HttpSettingsResp, err error) {
	data, err := setting.ParseAppHttpSettings()
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if err = copier.Copy(&resp, &data); err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
