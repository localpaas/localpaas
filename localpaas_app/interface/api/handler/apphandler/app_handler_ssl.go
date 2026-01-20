package apphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

// ObtainDomainSsl Obtains domain SSL
// @Summary Obtains domain SSL
// @Description Obtains domain SSL
// @Tags    apps
// @Produce json
// @Id      obtainDomainSsl
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   body body appdto.ObtainDomainSslReq true "request data"
// @Success 200 {object} appdto.ObtainDomainSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/ssl/obtain [post]
func (h *AppHandler) ObtainDomainSsl(ctx *gin.Context) {
	auth, projectID, appID, err := h.getAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewObtainDomainSslReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.ObtainDomainSsl(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
