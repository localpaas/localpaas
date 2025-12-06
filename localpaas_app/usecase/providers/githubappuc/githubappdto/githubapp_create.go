package githubappdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type CreateGithubAppReq struct {
	*GithubAppBaseReq
}

type GithubAppBaseReq struct {
	Name           string `json:"name"`
	ClientID       string `json:"clientId"`
	ClientSecret   string `json:"clientSecret"`
	Organization   string `json:"organization"`
	WebhookURL     string `json:"webhookURL"`
	WebhookSecret  string `json:"webhookSecret"`
	AppID          string `json:"appId"`
	InstallationID string `json:"installationId"`
	PrivateKey     string `json:"privateKey"`
	SSOEnabled     bool   `json:"ssoEnabled"`
}

func (req *GithubAppBaseReq) validate(_ string) []vld.Validator {
	// TODO: add validation
	return nil
}

func NewCreateGithubAppReq() *CreateGithubAppReq {
	return &CreateGithubAppReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateGithubAppReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateGithubAppResp struct {
	Meta *basedto.BaseMeta     `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}
