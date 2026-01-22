package oauthdto

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

type GetOAuthReq struct {
	settings.GetSettingReq
}

func NewGetOAuthReq() *GetOAuthReq {
	return &GetOAuthReq{}
}

func (req *GetOAuthReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetOAuthResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *OAuthResp        `json:"data"`
}

type OAuthResp struct {
	*settings.BaseSettingResp
	ClientID     string   `json:"clientId"`
	ClientSecret string   `json:"clientSecret"`
	Organization string   `json:"organization"`
	CallbackURL  string   `json:"callbackURL"`
	AuthURL      string   `json:"authURL,omitempty"`
	TokenURL     string   `json:"tokenURL,omitempty"`
	ProfileURL   string   `json:"profileURL,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
	Encrypted    bool     `json:"encrypted,omitempty"`
}

func (resp *OAuthResp) CopyClientSecret(field entity.EncryptedField) error {
	resp.ClientSecret = field.String()
	return nil
}

func TransformOAuth(setting *entity.Setting, baseCallbackURL string, objectID string) (resp *OAuthResp, err error) {
	config := setting.MustAsOAuth()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	// Recalculate callbackURL for the oauth as it depends on the actual server address
	resp.CallbackURL = baseCallbackURL + "/" + setting.ID
	resp.Encrypted = config.ClientSecret.IsEncrypted()
	if resp.Encrypted {
		resp.ClientSecret = maskedSecretKey
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
