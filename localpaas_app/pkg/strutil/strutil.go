package strutil

import (
	"github.com/samber/lo"
)

func ToSnakeCase(str string) string {
	return lo.SnakeCase(str)
}
