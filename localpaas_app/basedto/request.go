package basedto

import (
	"mime/multipart"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
)

// Paging is used to store pagination request from client side
type Paging struct {
	Sort   Orders
	Offset int
	Limit  int
}

func (p *Paging) Orders() Orders {
	return p.Sort
}

func (p *Paging) OffsetEnd() int {
	return p.Offset + p.Limit - 1
}

func (p *Paging) ToPage() int {
	return p.Offset / gofn.Coalesce(p.Limit, 1)
}

func (p *Paging) ToPageSize() int {
	return p.Limit
}

type Direction string

const (
	DirectionAsc  Direction = "ASC"
	DirectionDesc Direction = "DESC"
)

type Order struct {
	Direction  Direction
	ColumnName string
}

func (o *Order) OrderBy() string {
	return o.ColumnName + " " + string(o.Direction)
}

type Orders []*Order

func (o Orders) Contain(columnName string) bool {
	for i := range o {
		if o[i].ColumnName == columnName {
			return true
		}
	}
	return false
}

// Add merge given orders with skipping duplicated items
func (o *Orders) Add(orders ...*Order) {
	for _, order := range orders {
		if !o.Contain(order.ColumnName) {
			*o = append(*o, order)
		}
	}
}

func (o Orders) ApplyMapping(columnMapping map[string]string) {
	if len(columnMapping) == 0 {
		return
	}
	for _, order := range o {
		if mappedColumn, ok := columnMapping[order.ColumnName]; ok {
			if mappedColumn == "" {
				mappedColumn = order.ColumnName
			}
			order.ColumnName = mappedColumn
		}
	}
}

// ReqModifier any request struct implements this interface will have method `ModifyRequest()`
// to be called when request is parsed.
type ReqModifier interface {
	ModifyRequest() error
}

// ReqValidator any request struct implements this interface will have method `Validate()`
// to be called when request is parsed.
type ReqValidator interface {
	Validate() apperrors.ValidationErrors
}

// ReqParsingErrorHandler any request struct implements this interface will have method `HandleParsingError()`
// to be called when parsing error is raised.
type ReqParsingErrorHandler interface {
	HandleParsingError(err error) error
}

// ObjectIDReq request input for an object id
type ObjectIDReq struct {
	ID string `json:"id"`
}

func (req *ObjectIDReq) ToIDString() string {
	if req == nil {
		return ""
	}
	return req.ID
}

type ObjectIDSliceReq []*ObjectIDReq

func (req ObjectIDSliceReq) ToIDStringSlice() []string {
	result := make([]string, 0, len(req))
	for _, obj := range req {
		result = append(result, obj.ID)
	}
	return result
}

func (req ObjectIDSliceReq) HasID(id string) bool {
	for _, obj := range req {
		if obj.ID == id {
			return true
		}
	}
	return false
}

func (req *ObjectIDSliceReq) AppendID(id string) {
	*req = append(*req, &ObjectIDReq{ID: id})
}

// ObjectAccessReq request input for requesting access to an object
type ObjectAccessReq struct {
	ObjectIDReq
	Access base.AccessActions `json:"access"`
}

type ObjectAccessSliceReq []*ObjectAccessReq

func (req ObjectAccessSliceReq) ToIDStringSlice() []string {
	result := make([]string, 0, len(req))
	for _, obj := range req {
		result = append(result, obj.ID)
	}
	return result
}

// ModuleAccessReq request input for requesting access to a module
type ModuleAccessReq struct {
	ModuleIDReq
	Access base.AccessActions `json:"access"`
}

type ModuleIDReq struct {
	ID string `json:"id"` // could be module name
}

type ModuleAccessSliceReq []*ModuleAccessReq

type FileReq struct {
	File *multipart.FileHeader
}
