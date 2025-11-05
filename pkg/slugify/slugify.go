package slugify

import (
	"strings"

	"github.com/mozillazg/go-slugify"
)

func Slugify(s string) string {
	return slugify.Slugify(s)
}

func SlugifyEx(s string, replacements []string, limit int) string {
	slug := slugify.Slugify(s)
	i := 0
	for i < len(replacements)-1 {
		slug = strings.ReplaceAll(slug, replacements[i], replacements[i+1])
		i += 2
	}
	if len(slug) <= limit {
		return slug
	}
	return slug[:limit]
}
