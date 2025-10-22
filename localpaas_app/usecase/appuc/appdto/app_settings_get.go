package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
)

type GetAppSettingsReq struct {
	ProjectID string `json:"-"`
	AppID     string `json:"-"`
}

func NewGetAppSettingsReq() *GetAppSettingsReq {
	return &GetAppSettingsReq{}
}

func (req *GetAppSettingsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetAppSettingsResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *AppSettingsResp  `json:"data"`
}

type AppSettingsResp struct {
	Test string `json:"test"`
}

func TransformAppSettings(settings *entity.AppSettings) (resp *AppSettingsResp, err error) {
	if err = copier.Copy(&resp, &settings); err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
