package projectdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

const (
	minProjectNameLen = 1
	maxProjectNameLen = 50
)

func validateProjectName(name *string, field string) []vld.Validator {
	return basedto.ValidateStr(name, true, minProjectNameLen, maxProjectNameLen, field)
}
