package githubappdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

const (
	maskedSecretKey = "****************"
)

type GetGithubAppReq struct {
	settings.GetSettingReq
}

func NewGetGithubAppReq() *GetGithubAppReq {
	return &GetGithubAppReq{}
}

func (req *GetGithubAppReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetGithubAppResp struct {
	Meta *basedto.Meta  `json:"meta"`
	Data *GithubAppResp `json:"data"`
}

type GithubAppResp struct {
	*settings.BaseSettingResp
	ClientID       string `json:"clientId"`
	ClientSecret   string `json:"clientSecret"`
	Organization   string `json:"organization"`
	CallbackURL    string `json:"callbackURL"`
	WebhookURL     string `json:"webhookURL"`
	WebhookSecret  string `json:"webhookSecret"`
	AppID          int64  `json:"appId"`
	InstallationID int64  `json:"installationId"`
	PrivateKey     string `json:"privateKey"`
	SSOEnabled     bool   `json:"ssoEnabled"`
	Encrypted      bool   `json:"encrypted,omitempty"`
}

func (resp *GithubAppResp) CopyClientSecret(field entity.EncryptedField) error {
	resp.ClientSecret = field.String()
	return nil
}

func (resp *GithubAppResp) CopyPrivateKey(field entity.EncryptedField) error {
	resp.PrivateKey = field.String()
	return nil
}

func TransformGithubApp(setting *entity.Setting, baseCallbackURL string, objectID string) (
	resp *GithubAppResp, err error) {
	config := setting.MustAsGithubApp()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	// Recalculate callbackURL for the github-app as it depends on the actual server address
	resp.CallbackURL = baseCallbackURL + "/" + setting.ID
	resp.Encrypted = config.ClientSecret.IsEncrypted()
	if resp.Encrypted {
		resp.ClientSecret = maskedSecretKey
		resp.WebhookSecret = maskedSecretKey
		resp.PrivateKey = maskedSecretKey
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
