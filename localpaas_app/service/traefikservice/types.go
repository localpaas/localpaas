package traefikservice

import "github.com/localpaas/localpaas/localpaas_app/entity"

type AppConfigData struct {
	HttpSettings *entity.AppHttpSettings
	RefObjects   *entity.RefObjects
}
