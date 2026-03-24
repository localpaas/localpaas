package cloudstoragedto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/cloudprovideruc/cloudproviderdto"
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
	Provider     *cloudproviderdto.CloudProviderResp `json:"provider,omitempty"`
	S3           *CloudStorageS3Resp                 `json:"s3"`
	SecretMasked bool                                `json:"secretMasked,omitempty"`
}

type CloudStorageS3Resp struct {
	Region   string `json:"region"`
	Bucket   string `json:"bucket"`
	Endpoint string `json:"endpoint"`
}

func TransformCloudStorage(
	setting *entity.Setting,
	refObjects *entity.RefObjects,
) (resp *CloudStorageResp, err error) {
	config := setting.MustAsCloudStorage()
	if err = copier.Copy(&resp, &config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Provider, err = cloudproviderdto.TransformCloudProvider(refObjects.RefSettings[config.Provider.ID], refObjects)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.SecretMasked = resp.Provider.SecretMasked

	return resp, nil
}
