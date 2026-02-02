package basedto

import (
	vld "github.com/tiendc/go-validator"
)

func ValidateEnvVarsReq(envVars []*EnvVarReq, field string) (res []vld.Validator) {
	res = append(res, vld.SliceUniqueBy(envVars, func(t *EnvVarReq) string {
		return t.Key
	}).OnError(
		vld.SetField(field, nil),
		vld.SetCustomKey("ERR_VLD_VALUES_NON_UNIQUE"),
	))
	// TODO: use regex to validate allowed chars of env names
	return res
}
