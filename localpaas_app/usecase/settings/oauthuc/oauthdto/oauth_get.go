package oauthdto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
)

const (
	maskedSecretKey = "****************"
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
	ID           string   `json:"id"`
	Kind         string   `json:"kind,omitempty"`
	Name         string   `json:"name"`
	ClientID     string   `json:"clientId"`
	ClientSecret string   `json:"clientSecret"`
	Organization string   `json:"organization"`
	CallbackURL  string   `json:"callbackURL,omitempty"`
	AuthURL      string   `json:"authURL,omitempty"`
	TokenURL     string   `json:"tokenURL,omitempty"`
	ProfileURL   string   `json:"profileURL,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
	Encrypted    bool     `json:"encrypted,omitempty"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ExpireAt  *time.Time `json:"expireAt,omitempty" copy:",nilonzero"`
}

func TransformOAuth(setting *entity.Setting, baseCallbackURL string, decrypt bool) (resp *OAuthResp, err error) {
	if err = copier.Copy(&resp, &setting); err != nil {
		return nil, apperrors.Wrap(err)
	}

	config, err := setting.ParseOAuth(decrypt)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	// Recalculate callbackURL for the oauth as it depends on the actual server address
	resp.CallbackURL = baseCallbackURL + "/" + setting.Name
	resp.Encrypted = config.IsEncrypted()
	if resp.Encrypted {
		resp.ClientSecret = maskedSecretKey
	}
	return resp, nil
}
