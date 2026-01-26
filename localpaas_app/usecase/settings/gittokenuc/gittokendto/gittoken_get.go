package gittokendto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
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
	Meta *basedto.Meta `json:"meta"`
	Data *GitTokenResp `json:"data"`
}

type GitTokenResp struct {
	*settings.BaseSettingResp
	User      string `json:"user"`
	Token     string `json:"token"`
	BaseURL   string `json:"baseURL"`
	Encrypted bool   `json:"encrypted,omitempty"`
}

func (resp *GitTokenResp) CopyToken(field entity.EncryptedField) error {
	resp.Token = field.String()
	return nil
}

func TransformGitToken(setting *entity.Setting) (resp *GitTokenResp, err error) {
	config := setting.MustAsGitToken()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Encrypted = config.Token.IsEncrypted()
	if resp.Encrypted {
		resp.Token = maskedSecret
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
