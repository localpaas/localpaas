package cloudstoragedto

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

type CreateCloudStorageReq struct {
	settings.CreateSettingReq
	*CloudStorageBaseReq
}

type CloudStorageBaseReq struct {
	Name     string                `json:"name"`
	Kind     base.CloudStorageKind `json:"kind"`
	Provider basedto.ObjectIDReq   `json:"provider"`
	S3       *CloudStorageS3Req    `json:"s3"`
}

func (req *CloudStorageBaseReq) ToEntity() *entity.CloudStorage {
	return &entity.CloudStorage{
		Provider: entity.ObjectID{ID: req.Provider.ID},
		S3:       req.S3.ToEntity(),
	}
}

type CloudStorageS3Req struct {
	Region   string `json:"region"`
	Bucket   string `json:"bucket"`
	Endpoint string `json:"endpoint"`
}

func (req *CloudStorageS3Req) ToEntity() *entity.CloudStorageS3 {
	return &entity.CloudStorageS3{
		Region:   req.Region,
		Bucket:   req.Bucket,
		Endpoint: req.Endpoint,
	}
}

func (req *CloudStorageS3Req) validate(field string) (res []vld.Validator) {
	if req == nil {
		return nil
	}
	if field != "" {
		field += "."
	}
	res = append(res, basedto.ValidateStr(&req.Region, false, 1, maxKeyLen, field+"region")...)
	res = append(res, basedto.ValidateStr(&req.Bucket, false, 1, maxKeyLen, field+"bucket")...)
	res = append(res, basedto.ValidateStr(&req.Endpoint, false, 1, maxKeyLen, field+"endpoint")...)
	return res
}

func (req *CloudStorageBaseReq) validate(field string) (res []vld.Validator) {
	if field != "" {
		field += "."
	}
	res = append(res, basedto.ValidateStr(&req.Name, true, 1, base.SettingNameMaxLen, field+"name")...)
	res = append(res, basedto.ValidateStrIn(&req.Kind, true, base.AllCloudStorageKinds, field+"kind")...)
	res = append(res, basedto.ValidateObjectIDReq(&req.Provider, true, field+"provider")...)
	res = append(res, req.S3.validate("s3")...)
	return res
}

func NewCreateCloudStorageReq() *CreateCloudStorageReq {
	return &CreateCloudStorageReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateCloudStorageReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateCloudStorageResp struct {
	Meta *basedto.Meta         `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}
