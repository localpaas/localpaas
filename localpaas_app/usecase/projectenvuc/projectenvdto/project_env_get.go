package projectenvdto

import (
	"time"

	vld "github.com/tiendc/go-validator"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
)

type GetProjectEnvReq struct {
	ProjectID    string `json:"-"`
	ProjectEnvID string `json:"-"`
}

func NewGetProjectEnvReq() *GetProjectEnvReq {
	return &GetProjectEnvReq{}
}

func (req *GetProjectEnvReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.ProjectEnvID, true, "projectEnvId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetProjectEnvResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *ProjectEnvResp   `json:"data"`
}

type ProjectEnvResp struct {
	ID     string             `json:"id"`
	Name   string             `json:"name"`
	Status base.ProjectStatus `json:"status"`
	Apps   []*ProjectAppResp  `json:"apps"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ProjectAppResp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProjectEnvBaseResp struct {
	ID     string             `json:"id"`
	Name   string             `json:"name"`
	Status base.ProjectStatus `json:"status"`
}

func TransformProjectEnv(projectEnv *entity.ProjectEnv) (resp *ProjectEnvResp, err error) {
	if err = copier.Copy(&resp, &projectEnv); err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}

func TransformProjectEnvsBase(projectEnvs []*entity.ProjectEnv) []*ProjectEnvBaseResp {
	return gofn.MapSlice(projectEnvs, func(projectEnv *entity.ProjectEnv) *ProjectEnvBaseResp {
		return &ProjectEnvBaseResp{
			ID:     projectEnv.ID,
			Name:   projectEnv.Name,
			Status: projectEnv.Status,
		}
	})
}
