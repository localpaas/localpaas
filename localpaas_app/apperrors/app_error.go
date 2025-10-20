package apperrors

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	goerrors "github.com/go-errors/errors"
	"github.com/hashicorp/go-multierror"

	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/pkg/translation"
)

// AppError represents an error type to be used for any issue within the app.
// This error type is designed to be able to carry much extra information
// and ability to translate the error message into a specific language.
type AppError interface {
	error

	// WithCause sets cause of the error
	WithCause(cause error) AppError
	// WithParam sets a custom param (the param will be translated when build message)
	WithParam(k string, v any) AppError
	// WithParams sets custom params
	WithParams(map[string]any) AppError
	// WithNTParam sets a custom but non-translation param
	WithNTParam(k string, v any) AppError
	// WithMsgLog sets log message (used for debug purpose)
	WithMsgLog(format string, args ...any) AppError

	// Message builds representation message
	Message(lang translation.Lang) (msg string, transErr error)
	// Build builds error info for JSON API recommendation
	Build(lang translation.Lang) *ErrorInfo

	// UnwrapTilRoot unwraps til the root error
	UnwrapTilRoot() error
}

// appError implements AppError interface
type appError struct {
	err      error
	cause    error
	params   map[string]any
	ntParams map[string]any // non-translation params
	msgLog   string
}

// Error implements `error` interface
func (e *appError) Error() string {
	return e.err.Error()
}

func (e *appError) WithCause(cause error) AppError {
	e.cause = cause
	return e
}

func (e *appError) WithParam(k string, v any) AppError {
	e.params[k] = v
	return e
}

func (e *appError) WithParams(m map[string]any) AppError {
	for k, v := range m {
		e.params[k] = v
	}
	return e
}

func (e *appError) WithNTParam(k string, v any) AppError {
	e.ntParams[k] = v
	return e
}

func (e *appError) WithMsgLog(format string, args ...any) AppError {
	in := fmt.Sprintf(format, args...)
	if e.msgLog == "" {
		e.msgLog = in
		return e
	}
	e.msgLog = fmt.Sprintf("%s\n%s", e.msgLog, in)
	return e
}

// Build - builder (status, code, title, detail)
func (e *appError) Build(lang translation.Lang) *ErrorInfo {
	errInfo := &ErrorInfo{}

	errInfo.Status = e.getMappingStatus()
	if errInfo.Status == 0 {
		errInfo.Status = http.StatusInternalServerError
		errInfo.Code = ErrInternalServer.Error()
	} else {
		errInfo.Code = e.err.Error()
	}

	detail, transErr := e.Message(lang)
	if transErr != nil {
		// This is not error, just notify dev team about missing translation
		notifyTranslationMissing(transErr, lang)
	}

	errInfo.Title = http.StatusText(errInfo.Status)
	errInfo.Detail = detail
	errInfo.DebugLog = e.msgLog
	if e.cause != nil {
		errInfo.Cause = e.cause.Error()
	} else {
		errInfo.Cause = e.err.Error()
	}
	errInfo.StackTrace = e.StackTrace()

	return errInfo
}

func (e *appError) Message(lang translation.Lang) (msg string, transErr error) {
	params := make(map[string]any, len(e.params)+len(e.ntParams))
	for k, v := range e.ntParams {
		params[k] = v
	}
	for k, v := range e.params {
		vAsStr, ok := v.(string)
		if !ok {
			params[k] = v
			continue
		}
		if translated, err := translation.Localize(lang, vAsStr); err != nil {
			transErr = multierror.Append(transErr, err)
			params[k] = v
		} else {
			params[k] = translated
		}
	}

	// All of handled errors should have translation to guide users.
	// If it's the error that doesn't have translation, it's unexpected error.
	if !e.isErrHavingTranslation() {
		msg, _ := translation.Localize(lang, ErrInternalServer.Error())
		return msg, transErr //nolint:wrapcheck
	}

	msgID := e.UnwrapTilRoot().Error()
	if msg, err := translation.LocalizeEx(lang, msgID, params); err != nil {
		transErr = multierror.Append(transErr, err)
		// When the translation is missing, we cannot show correct message to users.
		// That's why just show 500 message. At least we can notify to define translations.
		msg, _ := translation.Localize(lang, ErrInternalServer.Error())
		return msg, transErr //nolint:wrapcheck
	} else {
		return msg, transErr //nolint:wrapcheck
	}
}

// Is - implements errors.Is.
// This returns true if either the inner error or the cause satisfies.
func (e *appError) Is(err error) bool {
	if errors.Is(e.err, err) {
		return true
	}
	if e.cause != nil {
		return errors.Is(e.cause, err)
	}
	return false
}

// Unwrap - implements errors.Unwrap
func (e *appError) Unwrap() error {
	return e.err
}

// UnwrapTilRoot keeps unwrapping until the root error
func (e *appError) UnwrapTilRoot() error {
	lastErr := e.err
	for {
		err := errors.Unwrap(lastErr)
		if err == nil {
			return lastErr
		}
		lastErr = err
	}
}

func (e *appError) StackTrace() string {
	var errWithStack *goerrors.Error
	if errors.As(e.err, &errWithStack) {
		return errWithStack.ErrorStack()
	}

	return ""
}

func (e *appError) getMappingStatus() int {
	// errorStatusMap[error] with an unhashable input error object
	// can cause panic. We recover from panic and return 0.
	defer func() {
		_ = recover()
	}()
	err := e.err
	for {
		if err == nil {
			return 0
		}
		if status, ok := errorStatusMap[err]; ok {
			return status
		}
		err = errors.Unwrap(err)
	}
}

func (e *appError) isErrHavingTranslation() bool {
	// Every error has prefix `ERR_` in regulation, and every defined error has translation.
	// Error key without having prefix `ERR_` won't be translated, they may be from a 3rd party library,
	// and we don't wrap them on purpose or a mistake.
	return strings.HasPrefix(e.UnwrapTilRoot().Error(), "ERR_")
}

func notifyTranslationMissing(e error, _ translation.Lang) {
	// the error format is something like this:
	// 1 error occurred:
	// \* message "ERR_BAD_REQUEST" not found in language "en"
	// It does have a line break and the actual error starts after '*',
	// so let's take the substring after '*' for the logging.
	_, errMsg, _ := strings.Cut(e.Error(), "* ")
	logging.Errorf("%s", errMsg)
}

func New(err error) AppError {
	if err == nil {
		return nil
	}
	var e *appError
	if errors.As(err, &e) {
		return e // already is a AppError, no need to wrap
	}
	return &appError{
		ntParams: map[string]any{},
		params:   map[string]any{},
		err:      goerrors.Wrap(err, 1),
	}
}
