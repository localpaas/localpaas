package httputil

import (
	"strings"

	"golang.org/x/text/language"

	"github.com/localpaas/localpaas/localpaas_app/pkg/translation"
)

// ParseRequestLang parse accept-language from http request.
// Return the best fit language for the input, if nothing matches the input, return the default language.
// E.g. Accept-Language: *
// E.g. Accept-Language: fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5
func ParseRequestLang(acceptLang string) translation.Lang {
	tags, _, err := language.ParseAcceptLanguage(acceptLang)
	if err != nil {
		return translation.GetDefaultLang()
	}
	for _, tag := range tags {
		lang := translation.Lang(tag.String())
		if ok := translation.IsAvailable(lang); ok {
			return lang
		}
	}
	return translation.GetDefaultLang()
}

// SplitQueryParam splits query param and removes all empty item
func SplitQueryParam[T ~string](str, sep string) (result []T) {
	for _, s := range strings.Split(str, sep) {
		s = strings.TrimSpace(s)
		if s != "" {
			result = append(result, T(s))
		}
	}
	return
}
