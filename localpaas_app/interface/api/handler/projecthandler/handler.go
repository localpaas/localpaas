package projecthandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/basesettinghandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/networkuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/volumeuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/projectuc"
)

type Handler struct {
	*basesettinghandler.Handler
	projectUC       *projectuc.UC
	dockerNetworkUC *networkuc.UC
	dockerVolumeUC  *volumeuc.UC
}

func New(
	baseSettingHandler *basesettinghandler.Handler,
	projectUC *projectuc.UC,
	dockerNetworkUC *networkuc.UC,
	dockerVolumeUC *volumeuc.UC,
) *Handler {
	return &Handler{
		Handler:         baseSettingHandler,
		projectUC:       projectUC,
		dockerNetworkUC: dockerNetworkUC,
		dockerVolumeUC:  dockerVolumeUC,
	}
}
