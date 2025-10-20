package entityutil

import (
	"strings"

	"github.com/localpaas/localpaas/localpaas_app/entity"
)

func SliceToIDMap[T entity.IDEntity, S ~[]T](entities S) map[string]T {
	result := make(map[string]T, len(entities))
	for _, ent := range entities {
		result[ent.GetID()] = ent
	}
	return result
}

func SliceToNameMap[T entity.NamedEntity, S ~[]T](entities S, caseSensitive bool) map[string]T {
	result := make(map[string]T, len(entities))
	for _, ent := range entities {
		if caseSensitive {
			result[ent.GetName()] = ent
		} else {
			result[strings.ToLower(ent.GetName())] = ent
		}
	}
	return result
}

// SliceToCaseSensitiveNameMap converts a slice of entities to a case-sensitive map of entities
// using the entity's name as a key.
func SliceToCaseSensitiveNameMap[T entity.NamedEntity, S ~[]T](entities S) map[string]T {
	return SliceToNameMap(entities, true)
}

// SliceToCaseISensitiveNameMap converts a slice of entities to a case-insensitive map of entities
// using the entity's name as a key.
func SliceToCaseISensitiveNameMap[T entity.NamedEntity, S ~[]T](entities S) map[string]T {
	return SliceToNameMap(entities, false)
}
