package awss3dto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

const (
	maxKeyLen = 100
)

type CreateAWSS3Req struct {
	settings.CreateSettingReq
	*AWSS3BaseReq
}

type AWSS3BaseReq struct {
	Name     string              `json:"name"`
	Cred     basedto.ObjectIDReq `json:"cred"`
	Region   string              `json:"region"`
	Bucket   string              `json:"bucket"`
	Endpoint string              `json:"endpoint"`
}

func (req *AWSS3BaseReq) ToEntity() *entity.AWSS3 {
	return &entity.AWSS3{
		Cred:     entity.ObjectID{ID: req.Cred.ID},
		Region:   req.Region,
		Bucket:   req.Bucket,
		Endpoint: req.Endpoint,
	}
}

func (req *AWSS3BaseReq) validate(field string) (res []vld.Validator) {
	if field != "" {
		field += "."
	}
	res = append(res, basedto.ValidateStr(&req.Name, true, 1, base.SettingNameMaxLen, field+"name")...)
	res = append(res, basedto.ValidateObjectIDReq(&req.Cred, true, field+"cred")...)
	res = append(res, basedto.ValidateStr(&req.Region, false, 1, maxKeyLen, field+"region")...)
	res = append(res, basedto.ValidateStr(&req.Bucket, false, 1, maxKeyLen, field+"bucket")...)
	res = append(res, basedto.ValidateStr(&req.Endpoint, false, 1, maxKeyLen, field+"endpoint")...)
	return res
}

func NewCreateAWSS3Req() *CreateAWSS3Req {
	return &CreateAWSS3Req{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateAWSS3Req) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateAWSS3Resp struct {
	Meta *basedto.Meta         `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}
