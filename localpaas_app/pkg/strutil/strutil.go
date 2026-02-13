package strutil

import (
	"github.com/samber/lo"
	"github.com/tiendc/gofn"
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

func CutShort(s string, maxLen int, padding string) string {
	if gofn.RuneLength(s) <= maxLen {
		return s
	}
	return string([]rune(s)[:maxLen]) + padding
}
