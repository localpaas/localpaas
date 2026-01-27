package emaildto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type CreateEmailReq struct {
	settings.CreateSettingReq
	*EmailBaseReq
}

type EmailBaseReq struct {
	Name string        `json:"name"`
	SMTP *SMTPConf     `json:"smtp"`
	HTTP *HTTPMailConf `json:"http"`
}

type SMTPConf struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
	SSL         bool   `json:"ssl"`
}

type HTTPMailConf struct {
	Endpoint    string            `json:"endpoint"`
	Method      string            `json:"method"`
	ContentType string            `json:"contentType"`
	Headers     map[string]string `json:"headers"`
	BodyMapping map[string]string `json:"bodyMapping"`
}

func (req *EmailBaseReq) validate(_ string) []vld.Validator {
	// TODO: add validation
	return nil
}

func NewCreateEmailReq() *CreateEmailReq {
	return &CreateEmailReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateEmailReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateEmailResp struct {
	Meta *basedto.Meta         `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}
