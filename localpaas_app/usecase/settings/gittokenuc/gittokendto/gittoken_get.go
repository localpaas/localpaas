package gittokendto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

const (
	maskedSecret = "****************"
)

type GetGitTokenReq struct {
	settings.GetSettingReq
}

func NewGetGitTokenReq() *GetGitTokenReq {
	return &GetGitTokenReq{}
}

func (req *GetGitTokenReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetGitTokenResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *GitTokenResp     `json:"data"`
}

type GitTokenResp struct {
	ID        string             `json:"id"`
	Kind      base.GitSource     `json:"kind"`
	Name      string             `json:"name"`
	Status    base.SettingStatus `json:"status"`
	User      string             `json:"user"`
	Token     string             `json:"token"`
	BaseURL   string             `json:"baseURL"`
	Encrypted bool               `json:"encrypted,omitempty"`
	UpdateVer int                `json:"updateVer"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ExpireAt  *time.Time `json:"expireAt,omitempty" copy:",nilonzero"`
}

func (resp *GitTokenResp) CopyToken(field entity.EncryptedField) error {
	resp.Token = field.String()
	return nil
}

func TransformGitToken(setting *entity.Setting) (resp *GitTokenResp, err error) {
	if err = copier.Copy(&resp, &setting); err != nil {
		return nil, apperrors.Wrap(err)
	}

	config := setting.MustAsGitToken()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Encrypted = config.Token.IsEncrypted()
	if resp.Encrypted {
		resp.Token = maskedSecret
	}
	return resp, nil
}
