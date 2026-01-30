package awss3dto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/pkg/entityutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/awsuc/awsdto"
)

type GetAWSS3Req struct {
	settings.GetSettingReq
}

func NewGetAWSS3Req() *GetAWSS3Req {
	return &GetAWSS3Req{}
}

func (req *GetAWSS3Req) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetAWSS3Resp struct {
	Meta *basedto.Meta `json:"meta"`
	Data *AWSS3Resp    `json:"data"`
}

type AWSS3Resp struct {
	*settings.BaseSettingResp
	Kind      string          `json:"kind,omitempty"`
	Cred      *awsdto.AWSResp `json:"cred,omitempty"`
	Region    string          `json:"region"`
	Bucket    string          `json:"bucket"`
	Endpoint  string          `json:"endpoint"`
	Encrypted bool            `json:"encrypted,omitempty"`
}

func TransformAWSS3(setting *entity.Setting) (resp *AWSS3Resp, err error) {
	config := setting.MustAsAWSS3()
	if err = copier.Copy(&resp, &config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	refAWS := entityutil.FindByID(setting.RefSettings, config.Cred.ID)
	if refAWS != nil {
		resp.Cred, err = awsdto.TransformAWS(refAWS)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
	}
	resp.Encrypted = resp.Cred.Encrypted

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
