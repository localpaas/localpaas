package ssldto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
)

const (
	maskedKey = "****************"
)

type GetSslReq struct {
	ID string `json:"-"`
}

func NewGetSslReq() *GetSslReq {
	return &GetSslReq{}
}

func (req *GetSslReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetSslResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *SslResp          `json:"data"`
}

type SslResp struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Certificate string     `json:"certificate"`
	PrivateKey  string     `json:"privateKey"`
	Expiration  *time.Time `json:"expiration" copy:",nilonzero"`
	Encrypted   bool       `json:"encrypted,omitempty"`
}

func TransformSsl(setting *entity.Setting, decrypt bool) (resp *SslResp, err error) {
	if err = copier.Copy(&resp, &setting); err != nil {
		return nil, apperrors.Wrap(err)
	}

	config, err := setting.ParseSsl(decrypt)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Encrypted = config.IsEncrypted()
	if resp.Encrypted {
		resp.PrivateKey = maskedKey
	}

	return resp, nil
}
