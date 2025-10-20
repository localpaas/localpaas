package apperrors

import (
	"errors"
	"net/http"
)

// Base errors with equivalent http status code
var (
	ErrInternalServer       = errors.New("ERR_INTERNAL_SERVER")
	ErrBadRequest           = errors.New("ERR_BAD_REQUEST")
	ErrParamInvalid         = errors.New("ERR_PARAM_INVALID")
	ErrUnavailable          = errors.New("ERR_UNAVAILABLE")
	ErrUnauthorized         = errors.New("ERR_UNAUTHORIZED")
	ErrForbidden            = errors.New("ERR_FORBIDDEN")
	ErrNotFound             = errors.New("ERR_NOT_FOUND")
	ErrAlreadyExist         = errors.New("ERR_ALREADY_EXIST")
	ErrConflict             = errors.New("ERR_CONFLICT")
	ErrNonEditable          = errors.New("ERR_NON_EDITABLE")
	ErrNonDeletable         = errors.New("ERR_NON_DELETABLE")
	ErrResourceInUse        = errors.New("ERR_RESOURCE_IN_USE")
	ErrRequestTooFrequently = errors.New("ERR_REQUEST_TOO_FREQUENTLY")
	ErrActionNotAllowed     = errors.New("ERR_ACTION_NOT_ALLOWED")
	ErrNotImplemented       = errors.New("ERR_NOT_IMPLEMENTED")
	ErrValidation           = errors.New("ERR_VALIDATION")
)

// Errors for session
var (
	ErrSessionJWTInvalid           = errors.New("ERR_SESSION_JWT_INVALID")
	ErrSessionJWTExpired           = errors.New("ERR_SESSION_JWT_EXPIRED")
	ErrSessionRefreshTokenRequired = errors.New("ERR_SESSION_REFRESH_TOKEN_REQUIRED")
	ErrSSORequired                 = errors.New("ERR_SSO_REQUIRED")
	ErrLoginInputInvalid           = errors.New("ERR_LOGIN_INPUT_INVALID")
	ErrPasswordMismatched          = errors.New("ERR_PASSWORD_MISMATCHED")
	ErrPasscodeMismatched          = errors.New("ERR_PASSCODE_MISMATCHED")
	ErrMFATokenInvalid             = errors.New("ERR_MFA_TOKEN_INVALID")
	ErrTooManyLoginFailures        = errors.New("ERR_TOO_MANY_LOGIN_FAILURES")
	ErrTooManyPasscodeAttempts     = errors.New("ERR_TOO_MANY_PASSCODE_ATTEMPTS")
)

// Errors for api client
var (
	ErrAPIKeyInactive   = errors.New("ERR_API_KEY_INACTIVE")
	ErrAPIKeyMismatched = errors.New("ERR_API_KEY_MISMATCHED")
	ErrAPIKeyInvalid    = errors.New("ERR_API_KEY_INVALID")
)

// Errors for user/workspace
var (
	ErrUserUnavailable                  = errors.New("ERR_USER_UNAVAILABLE")
	ErrUserStatusNotAllowSignUp         = errors.New("ERR_USER_STATUS_NOT_ALLOW_SIGN_UP")
	ErrUserStatusNotAllowActivation     = errors.New("ERR_USER_STATUS_NOT_ALLOW_ACTIVATION")
	ErrUserAlreadySignUp                = errors.New("ERR_USER_ALREADY_SIGN_UP")
	ErrUploadPhotoFailed                = errors.New("ERR_UPLOAD_PHOTO_FAILED")
	ErrPasswordNotMeetRequirements      = errors.New("ERR_PASSWORD_NOT_MEET_REQUIREMENTS")
	ErrPasswordResetTokenInvalid        = errors.New("ERR_PASSWORD_RESET_TOKEN_INVALID")
	ErrEmailChangeUnallowedOnEnforceSSO = errors.New("ERR_EMAIL_CHANGE_UNALLOWED_ON_ENFORCE_SSO")
	ErrEmailUnavailable                 = errors.New("ERR_EMAIL_UNAVAILABLE")
	ErrProjectAccessRequired            = errors.New("ERR_PROJECT_ACCESS_REQUIRED")
)

// errorStatusMap - mapping from error to http status code
var errorStatusMap = map[error]int{
	// Base errors
	ErrInternalServer:       http.StatusInternalServerError,
	ErrBadRequest:           http.StatusBadRequest,
	ErrParamInvalid:         http.StatusBadRequest,
	ErrUnavailable:          http.StatusBadRequest,
	ErrUnauthorized:         http.StatusUnauthorized,
	ErrForbidden:            http.StatusForbidden,
	ErrNotFound:             http.StatusNotFound,
	ErrAlreadyExist:         http.StatusConflict,
	ErrConflict:             http.StatusConflict,
	ErrNonEditable:          http.StatusUnprocessableEntity,
	ErrNonDeletable:         http.StatusUnprocessableEntity,
	ErrResourceInUse:        http.StatusConflict,
	ErrRequestTooFrequently: http.StatusForbidden,
	ErrActionNotAllowed:     http.StatusForbidden,
	ErrNotImplemented:       http.StatusNotImplemented,
	ErrValidation:           http.StatusBadRequest,

	// Session errors
	ErrSessionJWTInvalid:           http.StatusUnauthorized,
	ErrSessionJWTExpired:           http.StatusUnauthorized,
	ErrSessionRefreshTokenRequired: http.StatusForbidden,
	ErrSSORequired:                 http.StatusForbidden,
	ErrLoginInputInvalid:           http.StatusUnauthorized,
	ErrPasswordMismatched:          http.StatusUnauthorized,
	ErrPasscodeMismatched:          http.StatusUnauthorized,
	ErrMFATokenInvalid:             http.StatusUnauthorized,
	ErrTooManyLoginFailures:        http.StatusForbidden,
	ErrTooManyPasscodeAttempts:     http.StatusForbidden,

	// Api client errors
	ErrAPIKeyInactive:   http.StatusUnauthorized,
	ErrAPIKeyMismatched: http.StatusUnauthorized,
	ErrAPIKeyInvalid:    http.StatusUnauthorized,

	ErrUserUnavailable:                  http.StatusForbidden,
	ErrUserStatusNotAllowSignUp:         http.StatusForbidden,
	ErrUserStatusNotAllowActivation:     http.StatusForbidden,
	ErrUserAlreadySignUp:                http.StatusForbidden,
	ErrUploadPhotoFailed:                http.StatusUnprocessableEntity,
	ErrPasswordNotMeetRequirements:      http.StatusUnprocessableEntity,
	ErrPasswordResetTokenInvalid:        http.StatusNotAcceptable,
	ErrEmailChangeUnallowedOnEnforceSSO: http.StatusUnprocessableEntity,
	ErrEmailUnavailable:                 http.StatusUnprocessableEntity,
	ErrProjectAccessRequired:            http.StatusUnauthorized,
}

// errorWarnLevelMap defines the errors that are handled but unexpected to happen.
// Every error defined in this map will be notified at WARN level instead of ERROR.
var errorWarnLevelMap = map[error]bool{}
