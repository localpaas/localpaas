package emaildto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type CreateEmailReq struct {
	settings.CreateSettingReq
	*EmailBaseReq
}

type EmailBaseReq struct {
	Name string         `json:"name"`
	Kind base.EmailKind `json:"kind"`
	SMTP *EmailSMTP     `json:"smtp"`
	HTTP *EmailHTTP     `json:"http"`
}

func (req *EmailBaseReq) ToEntity() *entity.Email {
	email := &entity.Email{}
	switch req.Kind {
	case base.EmailKindSMTP:
		email.SMTP = req.SMTP.ToEntity()
	case base.EmailKindHTTP:
		email.HTTP = req.HTTP.ToEntity()
	}
	return email
}

type EmailSMTP struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
	SSL         bool   `json:"ssl"`
}

func (r *EmailSMTP) ToEntity() *entity.EmailSMTP {
	return &entity.EmailSMTP{
		Host:        r.Host,
		Port:        r.Port,
		Username:    r.Username,
		DisplayName: r.DisplayName,
		Password:    entity.NewEncryptedField(r.Password),
		SSL:         r.SSL,
	}
}

type EmailHTTP struct {
	Endpoint     string                        `json:"endpoint"`
	Method       string                        `json:"method"`
	ContentType  string                        `json:"contentType"`
	Headers      map[string]string             `json:"headers"`
	FieldMapping *entity.EmailHTTPFieldMapping `json:"fieldMapping"` // NOTE: use entity.EmailHTTPFieldMapping directly
	Username     string                        `json:"username"`
	DisplayName  string                        `json:"displayName"`
	Password     string                        `json:"password"`
}

func (r *EmailHTTP) ToEntity() *entity.EmailHTTP {
	return &entity.EmailHTTP{
		Endpoint:     r.Endpoint,
		Method:       r.Method,
		ContentType:  r.ContentType,
		Headers:      r.Headers,
		FieldMapping: r.FieldMapping,
		Username:     r.Username,
		DisplayName:  r.DisplayName,
		Password:     entity.NewEncryptedField(r.Password),
	}
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
