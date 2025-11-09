package s3storagedto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateS3StorageReq struct {
	ID string `json:"-"`
	*S3StoragePartialReq
}

type S3StoragePartialReq struct {
	Name            *string                      `json:"name"`
	AccessKeyID     *string                      `json:"accessKeyId"`
	SecretKey       *string                      `json:"secretKey"`
	Region          *string                      `json:"region"`
	Bucket          *string                      `json:"bucket"`
	ProjectAccesses []*S3StorageProjectAccessReq `json:"projectAccesses"`
}

func (req *S3StoragePartialReq) validate(field string) (res []vld.Validator) {
	if field != "" {
		field += "."
	}
	res = append(res, validateS3StorageName(req.Name, false, field+"name")...)
	res = append(res, basedto.ValidateStr(req.AccessKeyID, false, 1, maxKeyLen, "accessKeyId")...)
	res = append(res, basedto.ValidateStr(req.SecretKey, false, 1, maxKeyLen, "secretKey")...)
	res = append(res, basedto.ValidateStr(req.Region, false, 1, maxKeyLen, "region")...)
	res = append(res, basedto.ValidateStr(req.Bucket, false, 1, maxKeyLen, "bucket")...)
	return res
}

func NewUpdateS3StorageReq() *UpdateS3StorageReq {
	return &UpdateS3StorageReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateS3StorageReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateS3StorageResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}
