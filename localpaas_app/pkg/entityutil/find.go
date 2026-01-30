package entityutil

import (
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

func FindByID[T entity.IDEntity, S ~[]T](entities S, id string) (res T) {
	for _, ent := range entities {
		if ent.GetID() == id {
			return ent
		}
	}
	return
}

func FindByName[T entity.NamedEntity, S ~[]T](entities S, name string) (res T) {
	for _, ent := range entities {
		if ent.GetName() == name {
			return ent
		}
	}
	return
}
