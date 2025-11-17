package timeutil

import (
	"reflect"
	"time"

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

func MapstructureParseTimeFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}
		switch f.Kind() { //nolint:exhaustive
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string)) //nolint:forcetypeassert
		default:
			return data, nil
		}
	}
}

func MapstructureParseDurationFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if t != reflect.TypeOf(time.Duration(0)) {
			return data, nil
		}
		switch f.Kind() { //nolint:exhaustive
		case reflect.String:
			return time.ParseDuration(data.(string)) //nolint:forcetypeassert
		default:
			return data, nil
		}
	}
}
