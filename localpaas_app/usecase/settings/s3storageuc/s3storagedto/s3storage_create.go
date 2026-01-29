package s3storagedto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type CreateS3StorageReq struct {
	settings.CreateSettingReq
	*S3StorageBaseReq
}

type S3StorageBaseReq struct {
	Name        string `json:"name"`
	AccessKeyID string `json:"accessKeyId"`
	SecretKey   string `json:"secretKey"`
	Region      string `json:"region"`
	Bucket      string `json:"bucket"`
	Endpoint    string `json:"endpoint"`
}

func (req *S3StorageBaseReq) ToEntity() *entity.S3Storage {
	return &entity.S3Storage{
		AccessKeyID: req.AccessKeyID,
		SecretKey:   entity.NewEncryptedField(req.SecretKey),
		Region:      req.Region,
		Bucket:      req.Bucket,
		Endpoint:    req.Endpoint,
	}
}

func (req *S3StorageBaseReq) validate(field string) (res []vld.Validator) {
	if field != "" {
		field += "."
	}
	res = append(res, validateS3StorageName(&req.Name, true, field+"name")...)
	res = append(res, basedto.ValidateStr(&req.AccessKeyID, true, 1, maxKeyLen, "accessKeyId")...)
	res = append(res, basedto.ValidateStr(&req.SecretKey, true, 1, maxKeyLen, "secretKey")...)
	res = append(res, basedto.ValidateStr(&req.Region, false, 1, maxKeyLen, "region")...)
	res = append(res, basedto.ValidateStr(&req.Bucket, false, 1, maxKeyLen, "bucket")...)
	res = append(res, basedto.ValidateStr(&req.Endpoint, false, 1, maxKeyLen, "endpoint")...)
	return res
}

func NewCreateS3StorageReq() *CreateS3StorageReq {
	return &CreateS3StorageReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateS3StorageReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateS3StorageResp struct {
	Meta *basedto.Meta         `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}
