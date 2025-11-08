package oauthdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
)

type GetOAuthReq struct {
	ID string `json:"-"`
}

func NewGetOAuthReq() *GetOAuthReq {
	return &GetOAuthReq{}
}

func (req *GetOAuthReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetOAuthResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *OAuthResp        `json:"data"`
}

type OAuthResp struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Organization string `json:"organization"`
	BaseURL      string `json:"baseURL"`
	RedirectURL  string `json:"redirectURL"`
}

func TransformOAuth(setting *entity.Setting) (resp *OAuthResp, err error) {
	if err = copier.Copy(&resp, &setting); err != nil {
		return nil, apperrors.Wrap(err)
	}

	config, err := setting.ParseOAuth()
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	return resp, nil
}
