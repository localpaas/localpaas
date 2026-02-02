package basedto

import (
	"github.com/gitsight/go-vcsurl"
	vld "github.com/tiendc/go-validator"
)

func ValidateRepoURL[T ~string](s *T, required bool, field string) (result []vld.Validator) {
	if required {
		result = append(result, vld.Required(s).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_VALUE_REQUIRED"),
		))
	}
	if s != nil && *s != "" {
		_, err := vcsurl.Parse(string(*s))
		if err != nil {
			result = append(result, vld.Must(false).OnError(
				vld.SetField(field, nil),
				vld.SetCustomKey("ERR_VLD_URL_INVALID"),
			))
		}
	}
	return result
}
