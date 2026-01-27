package emaildto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

const (
	maskedPassword = "****************"
)

type GetEmailReq struct {
	settings.GetSettingReq
}

func NewGetEmailReq() *GetEmailReq {
	return &GetEmailReq{}
}

func (req *GetEmailReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetEmailResp struct {
	Meta *basedto.Meta `json:"meta"`
	Data *EmailResp    `json:"data"`
}

type EmailResp struct {
	*settings.BaseSettingResp
	Kind      base.EmailKind    `json:"kind"`
	SMTP      *SMTPConfResp     `json:"smtp,omitempty"`
	HTTP      *HTTPMailConfResp `json:"http,omitempty"`
	Encrypted bool              `json:"encrypted,omitempty"`
}

type SMTPConfResp struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
	SSL         bool   `json:"ssl"`
}

func (resp *SMTPConfResp) CopyPassword(field entity.EncryptedField) error {
	resp.Password = field.String()
	return nil
}

type HTTPMailConfResp struct {
	Endpoint    string            `json:"endpoint"`
	Method      string            `json:"method"`
	ContentType string            `json:"contentType"`
	Headers     map[string]string `json:"headers"`
	BodyMapping map[string]string `json:"bodyMapping"`
}

func TransformEmail(setting *entity.Setting) (resp *EmailResp, err error) {
	config := setting.MustAsEmail()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}
	resp.Kind = base.EmailKind(setting.Kind)

	switch {
	case config.SMTP != nil:
		resp.Encrypted = config.SMTP.Password.IsEncrypted()
		if resp.Encrypted {
			resp.SMTP.Password = maskedPassword
		}
	case config.HTTP != nil:
		// TODO: handle this?
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
