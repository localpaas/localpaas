package strutil

import (
	"strings"
	"unicode"

	"github.com/samber/lo"
)

func CapitalizeFirstLetter(s string) string {
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func ToSnakeCase(str string) string {
	return lo.SnakeCase(str)
}

func Quote(s string, quote string) string {
	return quote + s + quote
}

func Unquote(s string, quote string) string {
	return strings.TrimPrefix(strings.TrimSuffix(s, quote), quote)
}
