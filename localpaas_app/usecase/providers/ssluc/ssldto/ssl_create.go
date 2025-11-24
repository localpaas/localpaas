package ssldto

import (
	"strings"
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

const (
	nameMaxLen = 100
	keyMaxLen  = 10000
)

type CreateSslReq struct {
	*SslBaseReq
}

type SslBaseReq struct {
	Name        string           `json:"name"`
	Certificate string           `json:"certificate"`
	PrivateKey  string           `json:"privateKey"`
	KeySize     int              `json:"keySize"`
	Provider    base.SslProvider `json:"provider"`
	Email       string           `json:"email"`
	Expiration  time.Time        `json:"expiration"`
}

func (req *SslBaseReq) modifyRequest() error {
	req.Name = strings.TrimSpace(req.Name)
	return nil
}

func (req *SslBaseReq) validate(field string) (res []vld.Validator) {
	if field != "" {
		field += "."
	}
	res = append(res, basedto.ValidateStr(&req.Name, true, 1, nameMaxLen, field+"name")...)
	res = append(res, basedto.ValidateStr(&req.Certificate, true, 1, keyMaxLen, field+"certificate")...)
	res = append(res, basedto.ValidateStr(&req.PrivateKey, true, 1, keyMaxLen, field+"privateKey")...)
	res = append(res, basedto.ValidateStr(&req.Provider, false, 1, nameMaxLen, field+"provider")...)
	res = append(res, basedto.ValidateEmail(&req.Email, false, field+"email")...)
	return res
}

func NewCreateSslReq() *CreateSslReq {
	return &CreateSslReq{}
}

func (req *CreateSslReq) ModifyRequest() error {
	return req.modifyRequest()
}

// Validate implements interface basedto.ReqValidator
func (req *CreateSslReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateSslResp struct {
	Meta *basedto.BaseMeta     `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}
