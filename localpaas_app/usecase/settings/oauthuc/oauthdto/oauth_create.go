package oauthdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type CreateOAuthReq struct {
	OAuthType base.OAuthType `json:"oauthType"`
	*OAuthBaseReq
}

type OAuthBaseReq struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Organization string `json:"organization"`
	RedirectURL  string `json:"redirectURL"`
	BaseURL      string `json:"baseURL"`
}

func (req *OAuthBaseReq) validate(_ string) []vld.Validator {
	// TODO: add validation
	return nil
}

func NewCreateOAuthReq() *CreateOAuthReq {
	return &CreateOAuthReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateOAuthReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateStrIn(&req.OAuthType, true, base.AllOAuthTypes, "oauthType")...)
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateOAuthResp struct {
	Meta *basedto.BaseMeta     `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}
