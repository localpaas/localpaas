package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/cronjobuc/cronjobdto"
)

// ListCronJob Lists cron-jobs
// @Summary Lists cron-jobs
// @Description Lists cron-jobs
// @Tags    projects
// @Produce json
// @Id      listProjectCronJob
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} cronjobdto.ListCronJobResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/cron-jobs [get]
func (h *ProjectHandler) ListCronJob(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeCronJob, base.SettingScopeProject)
}

// GetCronJob Gets cron-job details
// @Summary Gets cron-job details
// @Description Gets cron-job details
// @Tags    projects
// @Produce json
// @Id      getProjectCronJob
// @Param   projectID path string true "project ID"
// @Param   id path string true "setting ID"
// @Success 200 {object} cronjobdto.GetCronJobResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/cron-jobs/{id} [get]
func (h *ProjectHandler) GetCronJob(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeCronJob, base.SettingScopeProject)
}

// CreateCronJob Creates a new cron-job
// @Summary Creates a new cron-job
// @Description Creates a new cron-job
// @Tags    projects
// @Produce json
// @Id      createProjectCronJob
// @Param   projectID path string true "project ID"
// @Param   body body cronjobdto.CreateCronJobReq true "request data"
// @Success 201 {object} cronjobdto.CreateCronJobResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/cron-jobs [post]
func (h *ProjectHandler) CreateCronJob(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeCronJob, base.SettingScopeProject)
}

// UpdateCronJob Updates cron-job
// @Summary Updates cron-job
// @Description Updates cron-job
// @Tags    projects
// @Produce json
// @Id      updateProjectCronJob
// @Param   projectID path string true "project ID"
// @Param   id path string true "setting ID"
// @Param   body body cronjobdto.UpdateCronJobReq true "request data"
// @Success 200 {object} cronjobdto.UpdateCronJobResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/cron-jobs/{id} [put]
func (h *ProjectHandler) UpdateCronJob(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeCronJob, base.SettingScopeProject)
}

// UpdateCronJobMeta Updates cron-job meta
// @Summary Updates cron-job meta
// @Description Updates cron-job meta
// @Tags    projects
// @Produce json
// @Id      updateProjectCronJobMeta
// @Param   projectID path string true "project ID"
// @Param   id path string true "setting ID"
// @Param   body body cronjobdto.UpdateCronJobMetaReq true "request data"
// @Success 200 {object} cronjobdto.UpdateCronJobMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/cron-jobs/{id}/meta [put]
func (h *ProjectHandler) UpdateCronJobMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeCronJob, base.SettingScopeProject)
}

// DeleteCronJob Deletes cron-job
// @Summary Deletes cron-job
// @Description Deletes cron-job
// @Tags    projects
// @Produce json
// @Id      deleteProjectCronJob
// @Param   projectID path string true "project ID"
// @Param   id path string true "setting ID"
// @Success 200 {object} cronjobdto.DeleteCronJobResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/cron-jobs/{id} [delete]
func (h *ProjectHandler) DeleteCronJob(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeCronJob, base.SettingScopeProject)
}
