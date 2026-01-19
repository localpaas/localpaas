package ssldto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
)

const (
	maskedKey = "****************"
)

type GetSslReq struct {
	providers.GetSettingReq
}

func NewGetSslReq() *GetSslReq {
	return &GetSslReq{}
}

func (req *GetSslReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetSslResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *SslResp          `json:"data"`
}

type SslResp struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Status      base.SettingStatus `json:"status"`
	Certificate string             `json:"certificate"`
	PrivateKey  string             `json:"privateKey"`
	KeySize     int                `json:"keySize"`
	Provider    string             `json:"provider"`
	Email       string             `json:"email"`
	Encrypted   bool               `json:"encrypted,omitempty"`
	UpdateVer   int                `json:"updateVer"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ExpireAt  *time.Time `json:"expireAt,omitempty" copy:",nilonzero"`
}

func (resp *SslResp) CopyPrivateKey(field entity.EncryptedField) error {
	resp.PrivateKey = field.String()
	return nil
}

func TransformSsl(setting *entity.Setting) (resp *SslResp, err error) {
	if err = copier.Copy(&resp, &setting); err != nil {
		return nil, apperrors.Wrap(err)
	}

	config := setting.MustAsSsl()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Encrypted = config.PrivateKey.IsEncrypted()
	if resp.Encrypted {
		resp.PrivateKey = maskedKey
	}
	return resp, nil
}
