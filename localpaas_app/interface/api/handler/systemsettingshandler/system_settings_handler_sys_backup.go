package systemsettingshandler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/fileuc/filedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/systemsettings/systembackupuc/systembackupdto"
)

// GetBackupSettings Gets backup settings
// @Summary Gets backup settings
// @Description Gets backup settings
// @Tags    system_settings
// @Produce json
// @Id      getSystemBackupSettings
// @Success 200 {object} systembackupdto.GetSystemBackupResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /system/settings/backup [get]
func (h *SystemSettingsHandler) GetBackupSettings(ctx *gin.Context) {
	auth, err := h.AuthHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSystem,
		ResourceType:   base.ResourceTypeSystemBackup,
		Action:         base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := systembackupdto.NewGetSystemBackupReq()
	req.Scope = base.NewSettingScopeGlobal()
	if err = h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.SystemBackupUC.GetSystemBackup(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateBackupSettings Updates backup settings
// @Summary Updates backup settings
// @Description Updates backup settings
// @Tags    system_settings
// @Produce json
// @Id      updateSystemBackupSettings
// @Param   body body systembackupdto.UpdateSystemBackupReq true "request data"
// @Success 200 {object} systembackupdto.UpdateSystemBackupResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /system/settings/backup [put]
func (h *SystemSettingsHandler) UpdateBackupSettings(ctx *gin.Context) {
	auth, err := h.AuthHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSystem,
		ResourceType:   base.ResourceTypeSystemBackup,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := systembackupdto.NewUpdateSystemBackupReq()
	req.Scope = base.NewSettingScopeGlobal()
	if err = h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.SystemBackupUC.UpdateSystemBackup(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// ExecuteBackup Executes the backup
// @Summary Executes the backup
// @Description Executes the backup
// @Tags    system_settings
// @Produce json
// @Id      executeSystemBackup
// @Param   body body systembackupdto.ExecuteSystemBackupReq true "request data"
// @Success 200 {object} systembackupdto.ExecuteSystemBackupResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /system/settings/backup/exec [post]
func (h *SystemSettingsHandler) ExecuteBackup(ctx *gin.Context) {
	auth, err := h.AuthHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSystem,
		ResourceType:   base.ResourceTypeSystemBackup,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := systembackupdto.NewExecuteSystemBackupReq()
	req.Scope = base.NewSettingScopeGlobal()
	if err = h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.SystemBackupUC.ExecuteSystemBackup(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// ListBackupFiles Lists backup files
// @Summary Lists backup files
// @Description Lists backup files
// @Tags    system_settings
// @Produce json
// @Id      listSystemBackupFiles
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} filedto.ListFileResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /system/settings/backup/files [get]
func (h *SystemSettingsHandler) ListBackupFiles(ctx *gin.Context) {
	auth, err := h.AuthHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSystem,
		ResourceType:   base.ResourceTypeSystemBackup,
		Action:         base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := filedto.NewListFileReq()
	req.Scope = base.NewSettingScopeGlobal()
	if err = h.ParseAndValidateRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.FileUC.ListFile(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetBackupFile Gets backup file
// @Summary Gets backup file
// @Description Gets backup file
// @Tags    system_settings
// @Produce json
// @Id      getSystemBackupFile
// @Param   fileID path string true "file setting ID"
// @Success 200 {object} filedto.GetFileResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /system/settings/backup/files/{fileID} [get]
func (h *SystemSettingsHandler) GetBackupFile(ctx *gin.Context) {
	fileID, err := h.ParseStringParam(ctx, "fileID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.AuthHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSystem,
		ResourceType:   base.ResourceTypeSystemBackup,
		ResourceID:     fileID,
		Action:         base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := filedto.NewGetFileReq()
	req.Scope = base.NewSettingScopeGlobal()
	req.ID = fileID
	if err = h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.FileUC.GetFile(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetBackupFileDownloadURL Gets download URL of a backup file
// @Summary Gets download URL of a backup file
// @Description Gets download URL of a backup file
// @Tags    system_settings
// @Produce json
// @Id      getSystemBackupFileDownloadURL
// @Param   fileID path string true "file setting ID"
// @Success 200 {object} filedto.GetFileDownloadURLResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /system/settings/backup/files/{fileID}/download-url [get]
func (h *SystemSettingsHandler) GetBackupFileDownloadURL(ctx *gin.Context) {
	fileID, err := h.ParseStringParam(ctx, "fileID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.AuthHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSystem,
		ResourceType:   base.ResourceTypeSystemBackup,
		ResourceID:     fileID,
		Action:         base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := filedto.NewGetFileDownloadURLReq()
	req.Scope = base.NewSettingScopeGlobal()
	req.ID = fileID
	req.RequireLogin = true
	req.CloudPresign = true
	req.Expiration = time.Minute * 5 //nolint:mnd
	if err = h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.FileUC.GetFileDownloadURL(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
