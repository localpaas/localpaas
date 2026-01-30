package awsdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type CreateAWSReq struct {
	settings.CreateSettingReq
	*AWSBaseReq
}

type AWSBaseReq struct {
	Name        string `json:"name"`
	AccessKeyID string `json:"accessKeyId"`
	SecretKey   string `json:"secretKey"`
	Region      string `json:"region"`
}

func (req *AWSBaseReq) ToEntity() *entity.AWS {
	return &entity.AWS{
		AccessKeyID: req.AccessKeyID,
		SecretKey:   entity.NewEncryptedField(req.SecretKey),
		Region:      req.Region,
	}
}

func (req *AWSBaseReq) validate(field string) (res []vld.Validator) {
	if field != "" {
		field += "."
	}
	res = append(res, validateAWSName(&req.Name, true, field+"name")...)
	res = append(res, basedto.ValidateStr(&req.AccessKeyID, true, 1, maxKeyLen, "accessKeyId")...)
	res = append(res, basedto.ValidateStr(&req.SecretKey, true, 1, maxKeyLen, "secretKey")...)
	res = append(res, basedto.ValidateStr(&req.Region, false, 1, maxKeyLen, "region")...)
	return res
}

func NewCreateAWSReq() *CreateAWSReq {
	return &CreateAWSReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateAWSReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateAWSResp struct {
	Meta *basedto.Meta         `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}
