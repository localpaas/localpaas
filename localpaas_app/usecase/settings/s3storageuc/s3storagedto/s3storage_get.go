package s3storagedto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

const (
	maskedSecretKey = "****************"
)

type GetS3StorageReq struct {
	settings.GetSettingReq
}

func NewGetS3StorageReq() *GetS3StorageReq {
	return &GetS3StorageReq{}
}

func (req *GetS3StorageReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetS3StorageResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *S3StorageResp    `json:"data"`
}

type S3StorageResp struct {
	ID          string             `json:"id"`
	Kind        string             `json:"kind,omitempty"`
	Name        string             `json:"name"`
	Status      base.SettingStatus `json:"status"`
	AccessKeyID string             `json:"accessKeyId"`
	SecretKey   string             `json:"secretKey,omitempty"`
	Region      string             `json:"region"`
	Bucket      string             `json:"bucket"`
	Endpoint    string             `json:"endpoint"`
	Encrypted   bool               `json:"encrypted,omitempty"`
	UpdateVer   int                `json:"updateVer"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ExpireAt  *time.Time `json:"expireAt,omitempty" copy:",nilonzero"`
}

func (resp *S3StorageResp) CopySecretKey(field entity.EncryptedField) error {
	resp.SecretKey = field.String()
	return nil
}

func TransformS3Storage(setting *entity.Setting) (resp *S3StorageResp, err error) {
	if err = copier.Copy(&resp, &setting); err != nil {
		return nil, apperrors.Wrap(err)
	}
	s3Config := setting.MustAsS3Storage()
	if err = copier.Copy(&resp, &s3Config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Encrypted = s3Config.SecretKey.IsEncrypted()
	if resp.Encrypted {
		resp.SecretKey = maskedSecretKey
	}
	return resp, nil
}
