package basedto

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/entity"
)

const (
	PageLimitDefault = 50
	PageLimitMax     = 10000

	CodeSuccess = "success"
)

type Response struct {
	Meta *Meta `json:"meta,omitempty"`
}

// BaseMeta metadata of single entity response
type BaseMeta struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// Meta metadata of response
type Meta struct {
	BaseMeta
	Page *PagingMeta `json:"page,omitempty"`
}

// PagingMeta metadata of pagination
type PagingMeta struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

// ObjectIDResp response for object with id
type ObjectIDResp struct {
	ID string `json:"id"`
}

type ObjectIDSliceResp []*ObjectIDResp

// NamedObjectResp response for object with name
type NamedObjectResp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ObjectUpdatedAtResp struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type NamedObjectSliceResp []*NamedObjectResp

type ObjectAccessResp struct {
	NamedObjectResp
	Access entity.AccessActions `json:"access"`
}

type ObjectAccessSliceResp []*ObjectAccessResp
