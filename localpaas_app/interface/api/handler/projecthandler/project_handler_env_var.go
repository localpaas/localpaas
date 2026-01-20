package projecthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc/projectdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// GetProjectEnvVars Gets project env vars
// @Summary Gets project env vars
// @Description Gets project env vars
// @Tags    projects
// @Produce json
// @Id      getProjectEnvVars
// @Param   projectID path string true "project ID"
// @Success 200 {object} projectdto.GetProjectEnvVarsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/env-vars [get]
func (h *ProjectHandler) GetProjectEnvVars(ctx *gin.Context) {
	auth, projectID, _, err := h.getAuth(ctx, base.ActionTypeRead, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := projectdto.NewGetProjectEnvVarsReq()
	req.ProjectID = projectID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.projectUC.GetProjectEnvVars(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateProjectEnvVars Updates project env vars
// @Summary Updates project env vars
// @Description Updates project env vars
// @Tags    projects
// @Produce json
// @Id      updateProjectEnvVars
// @Param   projectID path string true "project ID"
// @Param   body body projectdto.UpdateProjectEnvVarsReq true "request data"
// @Success 200 {object} projectdto.UpdateProjectEnvVarsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/env-vars [put]
func (h *ProjectHandler) UpdateProjectEnvVars(ctx *gin.Context) {
	auth, projectID, _, err := h.getAuth(ctx, base.ActionTypeWrite, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := projectdto.NewUpdateProjectEnvVarsReq()
	req.ProjectID = projectID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.projectUC.UpdateProjectEnvVars(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
