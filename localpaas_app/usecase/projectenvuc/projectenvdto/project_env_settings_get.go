package projectenvdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
)

type GetProjectEnvSettingsReq struct {
	ProjectID    string `json:"-"`
	ProjectEnvID string `json:"-"`
}

func NewGetProjectEnvSettingsReq() *GetProjectEnvSettingsReq {
	return &GetProjectEnvSettingsReq{}
}

func (req *GetProjectEnvSettingsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.ProjectEnvID, true, "projectEnvId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetProjectEnvSettingsResp struct {
	Meta *basedto.BaseMeta       `json:"meta"`
	Data *ProjectEnvSettingsResp `json:"data"`
}

type ProjectEnvSettingsResp struct {
	Test string `json:"test"`
}

func TransformProjectEnvSettings(settings *entity.ProjectEnvSettings) (resp *ProjectEnvSettingsResp, err error) {
	if err = copier.Copy(&resp, &settings); err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
