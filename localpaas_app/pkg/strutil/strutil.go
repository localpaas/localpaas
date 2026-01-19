package strutil

import (
	"github.com/samber/lo"
)

func ToSnakeCase(str string) string {
	return lo.SnakeCase(str)
}

func ToPascalCase(str string) string {
	return lo.PascalCase(str)
}

func ToCamelCase(str string) string {
	return lo.CamelCase(str)
}
