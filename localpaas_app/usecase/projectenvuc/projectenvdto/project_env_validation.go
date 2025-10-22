package projectenvdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

const (
	minProjectEnvNameLen = 1
	maxProjectEnvNameLen = 50
)

func validateProjectEnvName(name *string, field string) []vld.Validator {
	return basedto.ValidateStr(name, true, minProjectEnvNameLen, maxProjectEnvNameLen, field)
}
