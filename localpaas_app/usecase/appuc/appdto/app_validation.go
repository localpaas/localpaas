package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

const (
	minAppNameLen = 1
	maxAppNameLen = 100

	minAppTagLen = 0
	maxAppTagLen = 100

	minAppNoteLen = 10
	maxAppNoteLen = 10000
)

func validateAppName(name *string, field string) []vld.Validator {
	return basedto.ValidateStr(name, true, minAppNameLen, maxAppNameLen, field)
	// TODO: need validation for valid characters
}

func validateAppTags(tags []string, field string) []vld.Validator {
	return basedto.ValidateSliceEx(tags, true, minAppTagLen, maxAppTagLen, nil, field)
}

func validateAppNote(note *string, field string) []vld.Validator {
	return basedto.ValidateStr(note, false, minAppNoteLen, maxAppNoteLen, field)
}
