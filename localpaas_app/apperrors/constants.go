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
	ErrNameUnavailable      = errors.New("ERR_NAME_UNAVAILABLE")
	ErrUnauthorized         = errors.New("ERR_UNAUTHORIZED")
	ErrForbidden            = errors.New("ERR_FORBIDDEN")
	ErrNotFound             = errors.New("ERR_NOT_FOUND")
	ErrAlreadyExist         = errors.New("ERR_ALREADY_EXIST")
	ErrConflict             = errors.New("ERR_CONFLICT")
	ErrNonEditable          = errors.New("ERR_NON_EDITABLE")
	ErrNonDeletable         = errors.New("ERR_NON_DELETABLE")
	ErrResourceInUse        = errors.New("ERR_RESOURCE_IN_USE")
	ErrResourceInactive     = errors.New("ERR_RESOURCE_INACTIVE")
	ErrRequestTooFrequently = errors.New("ERR_REQUEST_TOO_FREQUENTLY")
	ErrActionNotAllowed     = errors.New("ERR_ACTION_NOT_ALLOWED")
	ErrNotImplemented       = errors.New("ERR_NOT_IMPLEMENTED")
	ErrUnsupported          = errors.New("ERR_UNSUPPORTED")
	ErrTokenInvalid         = errors.New("ERR_TOKEN_INVALID")
	ErrTypeInvalid          = errors.New("ERR_TYPE_INVALID")
	ErrValidation           = errors.New("ERR_VALIDATION")

	ErrUpdateVerMismatched = errors.New("ERR_UPDATE_VER_MISMATCHED")
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
	ErrTooManyLoginFailures        = errors.New("ERR_TOO_MANY_LOGIN_FAILURES")
	ErrTooManyPasscodeAttempts     = errors.New("ERR_TOO_MANY_PASSCODE_ATTEMPTS")
)

// Errors for api client
var (
	ErrAPIKeyMismatched = errors.New("ERR_API_KEY_MISMATCHED")
	ErrAPIKeyInvalid    = errors.New("ERR_API_KEY_INVALID")
)

// Errors for user
var (
	ErrUserUnavailable             = errors.New("ERR_USER_UNAVAILABLE")
	ErrUserStatusNotAllowAction    = errors.New("ERR_USER_STATUS_NOT_ALLOW_ACTION")
	ErrUserAlreadySignUp           = errors.New("ERR_USER_ALREADY_SIGN_UP")
	ErrUserNotCompleteMFASetup     = errors.New("ERR_USER_NOT_COMPLETE_MFA_SETUP")
	ErrPasswordNotMeetRequirements = errors.New("ERR_PASSWORD_NOT_MEET_REQUIREMENTS")
	ErrPasswordResetTokenInvalid   = errors.New("ERR_PASSWORD_RESET_TOKEN_INVALID")
	ErrEmailUnavailable            = errors.New("ERR_EMAIL_UNAVAILABLE")
	ErrEmailChangeUnallowed        = errors.New("ERR_EMAIL_CHANGE_UNALLOWED")
)

// nolint Errors from infrastructure
var (
	ErrInfra                   = errors.New("ERR_INFRA")
	ErrInfraUnknown            = errors.New("ERR_INFRA_UNKNOWN")
	ErrInfraInvalidArgument    = errors.Join(errors.New("ERR_INFRA_INVALID_ARGUMENT"), ErrParamInvalid)
	ErrInfraNotFound           = errors.Join(errors.New("ERR_INFRA_NOT_FOUND"), ErrNotFound)
	ErrInfraAlreadyExists      = errors.Join(errors.New("ERR_INFRA_ALREADY_EXISTS"), ErrAlreadyExist)
	ErrInfraPermissionDenied   = errors.New("ERR_INFRA_PERMISSION_DENIED")
	ErrInfraResourceExhausted  = errors.New("ERR_INFRA_RESOURCE_EXHAUSTED")
	ErrInfraFailedPrecondition = errors.New("ERR_INFRA_FAILED_PRECONDITION")
	ErrInfraConflict           = errors.Join(errors.New("ERR_INFRA_CONFLICT"), ErrConflict)
	ErrInfraNotModified        = errors.New("ERR_INFRA_NOT_MODIFIED")
	ErrInfraAborted            = errors.New("ERR_INFRA_ABORTED")
	ErrInfraOutOfRange         = errors.New("ERR_INFRA_OUT_OF_RANGE")
	ErrInfraNotImplemented     = errors.Join(errors.New("ERR_INFRA_NOT_IMPLEMENTED"), ErrNotImplemented)
	ErrInfraInternal           = errors.Join(errors.New("ERR_INFRA_INTERNAL"), ErrInternalServer)
	ErrInfraUnavailable        = errors.Join(errors.New("ERR_INFRA_UNAVAILABLE"), ErrUnavailable)
	ErrInfraDataLoss           = errors.New("ERR_INFRA_DATA_LOSS")
	ErrInfraUnauthorized       = errors.Join(errors.New("ERR_INFRA_UNAUTHORIZED"), ErrUnauthorized)
)

// Errors for cluster
var (
	ErrNodeRequiredByLocalPaasApp = errors.New("ERR_NODE_REQUIRED_BY_LOCAL_PAAS_APP")
)

// errorStatusMap - mapping from error to http status code
var errorStatusMap = map[error]int{
	// Base errors
	ErrInternalServer:       http.StatusInternalServerError,
	ErrBadRequest:           http.StatusBadRequest,
	ErrParamInvalid:         http.StatusBadRequest,
	ErrUnavailable:          http.StatusBadRequest,
	ErrUnauthorized:         http.StatusUnauthorized,
	ErrNameUnavailable:      http.StatusConflict,
	ErrForbidden:            http.StatusForbidden,
	ErrNotFound:             http.StatusNotFound,
	ErrAlreadyExist:         http.StatusConflict,
	ErrConflict:             http.StatusConflict,
	ErrNonEditable:          http.StatusUnprocessableEntity,
	ErrNonDeletable:         http.StatusUnprocessableEntity,
	ErrResourceInUse:        http.StatusConflict,
	ErrResourceInactive:     http.StatusNotAcceptable,
	ErrRequestTooFrequently: http.StatusForbidden,
	ErrActionNotAllowed:     http.StatusForbidden,
	ErrNotImplemented:       http.StatusNotImplemented,
	ErrUnsupported:          http.StatusNotImplemented,
	ErrTokenInvalid:         http.StatusUnauthorized,
	ErrTypeInvalid:          http.StatusInternalServerError,
	ErrValidation:           http.StatusBadRequest,
	ErrUpdateVerMismatched:  http.StatusConflict,

	// Session errors
	ErrSessionJWTInvalid:           http.StatusUnauthorized,
	ErrSessionJWTExpired:           http.StatusUnauthorized,
	ErrSessionRefreshTokenRequired: http.StatusForbidden,
	ErrSSORequired:                 http.StatusForbidden,
	ErrLoginInputInvalid:           http.StatusUnauthorized,
	ErrPasswordMismatched:          http.StatusUnauthorized,
	ErrPasscodeMismatched:          http.StatusUnauthorized,
	ErrTooManyLoginFailures:        http.StatusForbidden,
	ErrTooManyPasscodeAttempts:     http.StatusForbidden,

	// Api client errors
	ErrAPIKeyMismatched: http.StatusUnauthorized,
	ErrAPIKeyInvalid:    http.StatusUnauthorized,

	// User errors
	ErrUserUnavailable:             http.StatusForbidden,
	ErrUserStatusNotAllowAction:    http.StatusForbidden,
	ErrUserAlreadySignUp:           http.StatusForbidden,
	ErrUserNotCompleteMFASetup:     http.StatusForbidden,
	ErrPasswordNotMeetRequirements: http.StatusUnprocessableEntity,
	ErrPasswordResetTokenInvalid:   http.StatusNotAcceptable,
	ErrEmailUnavailable:            http.StatusUnprocessableEntity,
	ErrEmailChangeUnallowed:        http.StatusUnprocessableEntity,

	// Errors from infrastructure
	ErrInfra:                   http.StatusInternalServerError,
	ErrInfraUnknown:            http.StatusInternalServerError,
	ErrInfraInvalidArgument:    http.StatusBadRequest,
	ErrInfraNotFound:           http.StatusNotFound,
	ErrInfraAlreadyExists:      http.StatusConflict,
	ErrInfraPermissionDenied:   http.StatusUnauthorized,
	ErrInfraResourceExhausted:  http.StatusUnprocessableEntity,
	ErrInfraFailedPrecondition: http.StatusPreconditionFailed,
	ErrInfraConflict:           http.StatusConflict,
	ErrInfraNotModified:        http.StatusUnprocessableEntity,
	ErrInfraAborted:            http.StatusUnprocessableEntity,
	ErrInfraOutOfRange:         http.StatusUnprocessableEntity,
	ErrInfraNotImplemented:     http.StatusNotImplemented,
	ErrInfraInternal:           http.StatusInternalServerError,
	ErrInfraUnavailable:        http.StatusConflict,
	ErrInfraDataLoss:           http.StatusUnprocessableEntity,
	ErrInfraUnauthorized:       http.StatusUnauthorized,

	// Cluster errors
	ErrNodeRequiredByLocalPaasApp: http.StatusForbidden,
}

// errorWarnLevelMap defines the errors that are handled but unexpected to happen.
// Every error defined in this map will be notified at WARN level instead of ERROR.
var errorWarnLevelMap = map[error]bool{}

func RegisterStatusMapping(err error, statusCode int) {
	errorStatusMap[err] = statusCode
}
