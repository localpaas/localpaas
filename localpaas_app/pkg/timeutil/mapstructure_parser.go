package timeutil

import (
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
)

var (
	typeDate        = reflect.TypeFor[Date]()
	typeStdTime     = reflect.TypeFor[time.Time]()
	typeDuration    = reflect.TypeFor[Duration]()
	typeStdDuration = reflect.TypeFor[time.Duration]()

	// TODO: support pointer types such as *time.Time, etc
)

// MapstructureParseDateFunc custom decoder hook for mapstructure
//
//nolint:forcetypeassert
func MapstructureParseDateFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if t != typeDate {
			return data, nil
		}
		switch f.Kind() { //nolint:exhaustive
		case reflect.String:
			return ParseDate(data.(string))
		default:
			return data, nil
		}
	}
}

//nolint:forcetypeassert
func MapstructureParseTimeFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if t != typeStdTime {
			return data, nil
		}
		switch f.Kind() { //nolint:exhaustive
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		default:
			return data, nil
		}
	}
}

//nolint:forcetypeassert
func MapstructureParseDurationFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		switch f.Kind() { //nolint:exhaustive
		case reflect.String:
			if t == typeDuration {
				return ParseDuration(data.(string))
			}
			if t == typeStdDuration {
				return time.ParseDuration(data.(string))
			}
			return data, nil
		default:
			return data, nil
		}
	}
}
