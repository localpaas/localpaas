package cloudstoragedto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

const (
	maskedSecret = "********"
)

type GetCloudStorageReq struct {
	settings.GetSettingReq
}

func NewGetCloudStorageReq() *GetCloudStorageReq {
	return &GetCloudStorageReq{}
}

func (req *GetCloudStorageReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetCloudStorageResp struct {
	Meta *basedto.Meta     `json:"meta"`
	Data *CloudStorageResp `json:"data"`
}

type CloudStorageResp struct {
	*settings.BaseSettingResp
	S3           *CloudStorageS3Resp `json:"s3"`
	SecretMasked bool                `json:"secretMasked,omitempty"`
}

type CloudStorageS3Resp struct {
	*CloudProviderAWSResp
	Region   string `json:"region"`
	Bucket   string `json:"bucket"`
	Endpoint string `json:"endpoint"`
}

type CloudProviderAWSResp struct {
	AccessKeyID string `json:"accessKeyID"`
	SecretKey   string `json:"secretKey"`
	Region      string `json:"region"`
}

func TransformCloudStorage(
	setting *entity.Setting,
	_ *entity.RefObjects,
) (resp *CloudStorageResp, err error) {
	config := setting.MustAsCloudStorage()
	if err = copier.Copy(&resp, &config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.SecretMasked = resp.Inherited || (config.S3 != nil && config.S3.SecretKey.IsEncrypted())
	if resp.SecretMasked {
		if resp.S3 != nil {
			resp.S3.SecretKey = maskedSecret
		}
	}

	return resp, nil
}
