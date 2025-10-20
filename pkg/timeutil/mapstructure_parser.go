package timeutil

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// MapstructureParseDateFunc custom decoder hook for mapstructure
func MapstructureParseDateFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if t != reflect.TypeOf(Date{}) {
			return data, nil
		}
		switch f.Kind() { //nolint:exhaustive
		case reflect.String:
			return ParseDate(data.(string)) //nolint:forcetypeassert
		default:
			return data, nil
		}
	}
}
