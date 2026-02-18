package cacheentity

import (
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type HealthcheckSettings struct {
	Settings   []*entity.Setting  `json:"settings"`
	RefObjects *entity.RefObjects `json:"refObjects"`
}
