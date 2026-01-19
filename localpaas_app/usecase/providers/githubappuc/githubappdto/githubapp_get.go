package githubappdto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
)

const (
	maskedSecretKey = "****************"
)

type GetGithubAppReq struct {
	providers.GetSettingReq
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
	Meta *basedto.BaseMeta `json:"meta"`
	Data *GithubAppResp    `json:"data"`
}

type GithubAppResp struct {
	ID             string             `json:"id"`
	Kind           string             `json:"kind"`
	Name           string             `json:"name"`
	Status         base.SettingStatus `json:"status"`
	ClientID       string             `json:"clientId"`
	ClientSecret   string             `json:"clientSecret"`
	Organization   string             `json:"organization"`
	CallbackURL    string             `json:"callbackURL"`
	WebhookURL     string             `json:"webhookURL"`
	WebhookSecret  string             `json:"webhookSecret"`
	AppID          int64              `json:"appId"`
	InstallationID int64              `json:"installationId"`
	PrivateKey     string             `json:"privateKey"`
	SSOEnabled     bool               `json:"ssoEnabled"`
	Encrypted      bool               `json:"encrypted,omitempty"`
	UpdateVer      int                `json:"updateVer"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ExpireAt  *time.Time `json:"expireAt,omitempty" copy:",nilonzero"`
}

func (resp *GithubAppResp) CopyClientSecret(field entity.EncryptedField) error {
	resp.ClientSecret = field.String()
	return nil
}

func (resp *GithubAppResp) CopyWebhookSecret(field entity.EncryptedField) error {
	resp.WebhookSecret = field.String()
	return nil
}

func (resp *GithubAppResp) CopyPrivateKey(field entity.EncryptedField) error {
	resp.PrivateKey = field.String()
	return nil
}

func TransformGithubApp(setting *entity.Setting, baseCallbackURL string) (resp *GithubAppResp, err error) {
	if err = copier.Copy(&resp, &setting); err != nil {
		return nil, apperrors.Wrap(err)
	}

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
	return resp, nil
}
