package apppreviewhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/apppreviewuc/apppreviewdto"
)

// CreateAppPreview Creates preview for an app
// @Summary Creates preview for an app
// @Description Creates preview for an app
// @Tags    app_previews
// @Produce json
// @Id      createAppPreview
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   body body apppreviewdto.CreatePreviewReq true "request data"
// @Success 201 {object} apppreviewdto.CreatePreviewResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/previews [post]
func (h *Handler) CreateAppPreview(ctx *gin.Context) {
	auth, projectID, appID, _, err := h.GetAuthForItem(ctx, base.ActionTypeWrite, "")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := apppreviewdto.NewCreatePreviewReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appPreviewUC.CreatePreview(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}
