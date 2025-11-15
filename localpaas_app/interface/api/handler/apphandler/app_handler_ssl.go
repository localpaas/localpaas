package apphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// InstallDomainSsl Installs domain SSL
// @Summary Installs domain SSL
// @Description Installs domain SSL
// @Tags    apps
// @Produce json
// @Id      installDomainSsl
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   body body appdto.InstallDomainSslReq true "request data"
// @Success 200 {object} appdto.InstallDomainSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/ssl/install [post]
func (h *AppHandler) InstallDomainSsl(ctx *gin.Context) {
	projectID, err := h.ParseStringParam(ctx, "projectID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	appID, err := h.ParseStringParam(ctx, "appID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule:     base.ResourceModuleProject,
		ParentResourceType: base.ResourceTypeProject,
		ParentResourceID:   projectID,
		ResourceType:       base.ResourceTypeApp,
		ResourceID:         appID,
		Action:             base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewInstallDomainSslReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err := h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.InstallDomainSsl(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
