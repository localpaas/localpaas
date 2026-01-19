package oauthdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
)

type CreateOAuthReq struct {
	providers.CreateSettingReq
	*OAuthBaseReq
}

type OAuthBaseReq struct {
	Kind         base.OAuthKind `json:"kind"`
	Name         string         `json:"name"`
	ClientID     string         `json:"clientId"`
	ClientSecret string         `json:"clientSecret"`
	Organization string         `json:"organization"`
	AuthURL      string         `json:"authURL"`
	TokenURL     string         `json:"tokenURL"`
	ProfileURL   string         `json:"profileURL"`
	Scopes       []string       `json:"scopes"`
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
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateOAuthResp struct {
	Meta *basedto.BaseMeta  `json:"meta"`
	Data *OAuthCreationResp `json:"data"`
}

type OAuthCreationResp struct {
	ID          string `json:"id"`
	CallbackURL string `json:"callbackURL"`
}
