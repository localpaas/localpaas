package accessiblebyprojectsdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc/projectdto"
)

type GetAccessibleByProjectsReq struct {
	SettingID string `json:"-"`
}

func NewGetAccessibleByProjectsReq() *GetAccessibleByProjectsReq {
	return &GetAccessibleByProjectsReq{}
}

func (req *GetAccessibleByProjectsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.SettingID, true, "settingId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetAccessibleByProjectsResp struct {
	Meta *basedto.Meta                 `json:"meta"`
	Data *AccessibleByProjectsDataResp `json:"data"`
}

type AccessibleByProjectsDataResp struct {
	AccessibleByProjects []*AccessibleByProjectResp `json:"accessibleByProjects"`
}

type AccessibleByProjectResp struct {
	*projectdto.ProjectBaseResp
}

func TransformAccessibleByProjects(setting *entity.Setting) *AccessibleByProjectsDataResp {
	resp := &AccessibleByProjectsDataResp{
		AccessibleByProjects: make([]*AccessibleByProjectResp, 0, len(setting.AccessibleByProjects)),
	}
	for _, ap := range setting.AccessibleByProjects {
		if ap.Project == nil {
			continue
		}
		resp.AccessibleByProjects = append(resp.AccessibleByProjects,
			TransformAccessibleByProject(ap.Project))
	}
	return resp
}

func TransformAccessibleByProject(project *entity.Project) *AccessibleByProjectResp {
	return &AccessibleByProjectResp{
		ProjectBaseResp: projectdto.TransformProjectBase(project),
	}
}
