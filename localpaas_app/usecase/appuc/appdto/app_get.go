package appdto

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

type GetAppReq struct {
	ProjectID string `json:"-"`
	AppID     string `json:"-"`
}

func NewGetAppReq() *GetAppReq {
	return &GetAppReq{}
}

func (req *GetAppReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetAppResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *AppResp          `json:"data"`
}

type AppResp struct {
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	Key          string               `json:"key"`
	Status       base.AppStatus       `json:"status"`
	Photo        string               `json:"photo"`
	Note         string               `json:"note"`
	Tags         []string             `json:"tags" copy:"-"` // manual copy AppTag -> string
	UserAccesses []*AppUserAccessResp `json:"userAccesses"`
	UpdateVer    int                  `json:"updateVer"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AppUserAccessResp struct {
	*basedto.UserBaseResp
	Access base.AccessActions `json:"access"`
}

type AppBaseResp struct {
	ID     string         `json:"id"`
	Name   string         `json:"name"`
	Key    string         `json:"key"`
	Photo  string         `json:"photo"`
	Status base.AppStatus `json:"status"`
}

func TransformApp(app *entity.App) (resp *AppResp, err error) {
	if err = copier.Copy(&resp, &app); err != nil {
		return nil, apperrors.Wrap(err)
	}
	resp.Tags = gofn.MapSlice(app.Tags, func(t *entity.AppTag) string { return t.Tag })
	resp.UserAccesses = TransformUserAccesses(app.Accesses)
	return resp, nil
}

func TransformUserAccesses(accesses []*entity.ACLPermission) []*AppUserAccessResp {
	return gofn.MapSlice(accesses, TransformUserAccess)
}

func TransformUserAccess(access *entity.ACLPermission) *AppUserAccessResp {
	return &AppUserAccessResp{
		UserBaseResp: basedto.TransformUserBase(access.SubjectUser),
		Access:       access.Actions,
	}
}

func TransformAppsBase(apps []*entity.App) []*AppBaseResp {
	return gofn.MapSlice(apps, func(app *entity.App) *AppBaseResp {
		return &AppBaseResp{
			ID:     app.ID,
			Name:   app.Name,
			Key:    app.Key,
			Photo:  app.Photo,
			Status: app.Status,
		}
	})
}
