package projectdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
)

type GetProjectSettingsReq struct {
	ProjectID string `json:"-"`
}

func NewGetProjectSettingsReq() *GetProjectSettingsReq {
	return &GetProjectSettingsReq{}
}

func (req *GetProjectSettingsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetProjectSettingsResp struct {
	Meta *basedto.BaseMeta    `json:"meta"`
	Data *ProjectSettingsResp `json:"data"`
}

type ProjectSettingsResp struct {
	Test string `json:"test"`
}

func TransformProjectSettings(settings *entity.ProjectSettings) (resp *ProjectSettingsResp, err error) {
	if err = copier.Copy(&resp, &settings); err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}
