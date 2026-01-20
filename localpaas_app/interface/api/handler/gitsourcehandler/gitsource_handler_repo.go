package gitsourcehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/gitsourceuc/gitsourcedto"
)

// ListGitRepo Lists repo
// @Summary Lists repo
// @Description Lists repo
// @Tags    git_source
// @Produce json
// @Id      listGitRepo
// @Param   settingID path string true "github-app ID or git-token ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} gitsourcedto.ListRepoResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /git-source/{settingID}/repositories [get]
func (h *GitSourceHandler) ListGitRepo(ctx *gin.Context) {
	settingID, err := h.ParseStringParam(ctx, "settingID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSettings,
		ResourceType:   base.ResourceTypeGitToken,
		Action:         base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := gitsourcedto.NewListRepoReq()
	req.SettingID = settingID
	if err := h.ParseAndValidateRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.gitSourceUC.ListRepo(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
