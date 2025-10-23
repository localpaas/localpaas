package sshkeydto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

const (
	minNameLen = 1
	maxNameLen = 100

	maxKeyLen = 1000
)

func validateSSHKeyName(name *string, required bool, field string) []vld.Validator {
	if !required && (name == nil || *name == "") {
		return nil
	}
	return basedto.ValidateStr(name, true, minNameLen, maxNameLen, field)
	// TODO: need validation for valid characters
}
