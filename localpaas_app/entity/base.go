package entity

const (
	defaultSaltLen = 10
)

// IDEntity base interface for every entity kind which has `ID`
type IDEntity interface {
	GetID() string
}

// NamedEntity base interface for every entity kind which has `name`
type NamedEntity interface {
	GetName() string
}
