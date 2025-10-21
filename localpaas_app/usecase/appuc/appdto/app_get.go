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
	ID string `json:"-"`
}

func NewGetAppReq() *GetAppReq {
	return &GetAppReq{}
}

func (req *GetAppReq) Validate() apperrors.ValidationErrors {
	return apperrors.NewValidationErrors(vld.Validate(basedto.ValidateID(&req.ID, true, "id")...))
}

type GetAppResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *AppResp          `json:"data"`
}

type AppResp struct {
	ID     string         `json:"id"`
	Name   string         `json:"name"`
	Status base.AppStatus `json:"status"`
	Photo  string         `json:"photo"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AppBaseResp struct {
	ID     string         `json:"id"`
	Name   string         `json:"name"`
	Photo  string         `json:"photo"`
	Status base.AppStatus `json:"status"`
}

func TransformApp(app *entity.App) (resp *AppResp, err error) {
	if err = copier.Copy(&resp, &app); err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}

func TransformAppsBase(apps []*entity.App) []*AppBaseResp {
	return gofn.MapSlice(apps, func(app *entity.App) *AppBaseResp {
		return &AppBaseResp{
			ID:     app.ID,
			Name:   app.Name,
			Photo:  app.Photo,
			Status: app.Status,
		}
	})
}
