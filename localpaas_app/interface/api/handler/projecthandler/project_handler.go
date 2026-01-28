package projecthandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/basesettinghandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc"
)

type ProjectHandler struct {
	*basesettinghandler.BaseSettingHandler
	projectUC *projectuc.ProjectUC
}

func NewProjectHandler(
	baseSettingHandler *basesettinghandler.BaseSettingHandler,
	projectUC *projectuc.ProjectUC,
) *ProjectHandler {
	return &ProjectHandler{
		BaseSettingHandler: baseSettingHandler,
		projectUC:          projectUC,
	}
}
