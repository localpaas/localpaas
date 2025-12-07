package apperrors

import (
	"errors"
	"net/http"

	"github.com/localpaas/localpaas/localpaas_app/pkg/translation"
)

// Wrap wraps an error with storing stack trace
func Wrap(err error) error {
	return New(err)
}

type ErrLevel uint8

const (
	ErrLevelInfo  ErrLevel = iota + 1
	ErrLevelWarn  ErrLevel = iota + 1
	ErrLevelError ErrLevel = iota + 1
)

// ParseError parse the given error and return a list of ErrorInfo.
// If the given error is a single one, the returned slice will contain only one item.
func ParseError(err error, lang translation.Lang) (*ErrorInfo, ErrLevel) {
	// err is ValidationErrors
	validationErrs := ValidationErrors{}
	if errors.As(err, &validationErrs) {
		return validationErrs.Build(lang), ErrLevelInfo
	}

	// `New` will automatically create AppError if the input is not AppError
	appErr := New(err)
	errorInfo := appErr.Build(lang)
	if errorInfo.Status == http.StatusInternalServerError {
		return errorInfo, ErrLevelError
	}
	if errorWarnLevelMap[appErr.UnwrapTilRoot()] {
		return errorInfo, ErrLevelWarn
	}
	// User error, not the logic and not unexpected, reports at INFO level
	return errorInfo, ErrLevelInfo
}

// ParseErrorDetail parses to get detail from the given error
func ParseErrorDetail(err error, lang translation.Lang) (detail string) {
	errInfo, _ := ParseError(err, lang)
	if errInfo != nil {
		detail = errInfo.Detail
	}
	return detail
}

// NewNotFound return AppError for error NotFound
func NewNotFound(name string) AppError {
	return New(ErrNotFound).WithParam("Name", name)
}
func NewNotFoundNT(name string) AppError { // NT: non translation param
	return New(ErrNotFound).WithNTParam("Name", name)
}

// NewAlreadyExist return AppError for error AlreadyExist
func NewAlreadyExist(name string) AppError {
	return New(ErrAlreadyExist).WithParam("Name", name)
}
func NewAlreadyExistNT(name string) AppError { // NT: non translation param
	return New(ErrAlreadyExist).WithNTParam("Name", name)
}

// NewConflict return AppError for error Conflict
func NewConflict(name string) AppError {
	return New(ErrConflict).WithParam("Name", name)
}
func NewConflictNT(name string) AppError { // NT: non translation param
	return New(ErrConflict).WithNTParam("Name", name)
}

// NewParamInvalid return AppError for error ParamInvalid
func NewParamInvalid(name string) AppError {
	return New(ErrParamInvalid).WithParam("Name", name)
}
func NewParamInvalidNT(name string) AppError { // NT: non translation param
	return New(ErrParamInvalid).WithNTParam("Name", name)
}

// NewUnavailable return AppError for error Unavailable
func NewUnavailable(name string) AppError {
	return New(ErrUnavailable).WithParam("Name", name)
}
func NewUnavailableNT(name string) AppError { // NT: non translation param
	return New(ErrUnavailable).WithNTParam("Name", name)
}

// NewNonEditable return AppError for error NonEditable
func NewNonEditable(name string) AppError {
	return New(ErrNonEditable).WithParam("Name", name)
}
func NewNonEditableNT(name string) AppError { // NT: non translation param
	return New(ErrNonEditable).WithNTParam("Name", name)
}

// NewNonDeletable return AppError for error NonDeletable
func NewNonDeletable(name string) AppError {
	return New(ErrNonDeletable).WithParam("Name", name)
}
func NewNonDeletableNT(name string) AppError { // NT: non translation param
	return New(ErrNonDeletable).WithNTParam("Name", name)
}

// NewInUse return AppError for error ResourceInUse
func NewInUse(name string) AppError {
	return New(ErrResourceInUse).WithParam("Name", name)
}
func NewInUseNT(name string) AppError { // NT: non translation param
	return New(ErrResourceInUse).WithNTParam("Name", name)
}

// NewTypeInvalid return AppError for error TypeInvalid
func NewTypeInvalid() AppError {
	return New(ErrTypeInvalid)
}
func NewTypeInvalidNT() AppError { // NT: non translation param
	return New(ErrTypeInvalid)
}

// NewUnsupported return AppError for error Unsupported
func NewUnsupported() AppError {
	return New(ErrUnsupported)
}
func NewUnsupportedNT() AppError { // NT: non translation param
	return New(ErrUnsupported)
}

// NewNotImplemented return AppError for error NotImplemented
func NewNotImplemented() AppError {
	return New(ErrNotImplemented)
}
func NewNotImplementedNT() AppError { // NT: non translation param
	return New(ErrNotImplemented)
}
