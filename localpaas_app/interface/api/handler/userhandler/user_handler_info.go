package userhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/useruc/userdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListUserSimple Lists users
// @Summary Lists users
// @Description Lists users
// @Tags    users_info
// @Produce json
// @Id      listUserSimple
// @Param   filter[search] query string false "`filter[search]=full name (support *)` to search specific users"
// @Param   page[offset] query int false "`page[offset]=offset`"
// @Param   page[limit] query int false "`page[limit]=limit`"
// @Param   sort query string false "`sort=[-]firstName|lastName|createdAt|updatedAt|...`"
// @Success 200 {object} userdto.ListUserSimpleResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /users/base-list [get]
func (h *UserHandler) ListUserSimple(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceType: base.ResourceTypeUser,
		Action:       base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := userdto.NewListUserSimpleReq()
	if err = h.ParseRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.userUC.ListUserSimple(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetUser Gets user details
// @Summary Gets user details
// @Description Gets user details
// @Tags    users_info
// @Produce json
// @Id      getUser
// @Param   userID path string true "user ID"
// @Success 200 {object} userdto.GetUserResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /users/{userID} [get]
func (h *UserHandler) GetUser(ctx *gin.Context) {
	userID, err := h.ParseStringParam(ctx, "userID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceType: base.ResourceTypeUser,
		ResourceID:   userID,
		Action:       base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := userdto.NewGetUserReq()
	req.ID = userID
	if err = h.ParseRequest(ctx, req, nil); err != nil { // to make sure Validate() to be called
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.userUC.GetUser(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// ListUser Lists users of the current workspace
// @Summary Lists users of the current workspace
// @Description Lists users of the current workspace
// @Tags    users_info
// @Produce json
// @Id      listUser
// @Param   filter[roleId] query string false "`filter[roleId]=id1,id2` to filter by role"
// @Param   filter[type] query string false "`filter[type]=internal|external` to filter by type"
// @Param   filter[status] query string false "`filter[status]=active,invited,disabled` to filter by status"
// @Param   filter[search] query string false "`filter[search]=email|name|phone (support *)` to search specific users"
// @Param   page[offset] query int false "`page[offset]=offset`"
// @Param   page[limit] query int false "`page[limit]=limit`"
// @Param   sort query string false "`sort=[-]type|firstName|lastName|createdAt|updatedAt|...`"
// @Param   filters query string false "advanced filters in format `and(field.op(arg1, arg2), or(field.op(arg1))`"
// @Success 200 {object} userdto.ListUserResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /users [get]
func (h *UserHandler) ListUser(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceType: base.ResourceTypeUser,
		Action:       base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := userdto.NewListUserReq()
	if err = h.ParseRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.userUC.ListUser(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
