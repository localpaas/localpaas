package apperrors

import (
	"github.com/containerd/errdefs"
)

func NewInfra(err error) AppError {
	if err == nil {
		return nil
	}
	switch {
	case errdefs.IsUnknown(err):
		return New(ErrInfraUnknown).WithNTParam("Error", err.Error())
	case errdefs.IsInvalidArgument(err):
		return New(ErrInfraInvalidArgument).WithNTParam("Error", err.Error())
	case errdefs.IsNotFound(err):
		return New(ErrInfraNotFound).WithNTParam("Error", err.Error())
	case errdefs.IsAlreadyExists(err):
		return New(ErrInfraAlreadyExists).WithNTParam("Error", err.Error())
	case errdefs.IsPermissionDenied(err):
		return New(ErrInfraPermissionDenied).WithNTParam("Error", err.Error())
	case errdefs.IsResourceExhausted(err):
		return New(ErrInfraResourceExhausted).WithNTParam("Error", err.Error())
	case errdefs.IsFailedPrecondition(err):
		return New(ErrInfraFailedPrecondition).WithNTParam("Error", err.Error())
	case errdefs.IsConflict(err):
		return New(ErrInfraConflict).WithNTParam("Error", err.Error())
	case errdefs.IsNotModified(err):
		return New(ErrInfraNotModified).WithNTParam("Error", err.Error())
	case errdefs.IsAborted(err):
		return New(ErrInfraAborted).WithNTParam("Error", err.Error())
	case errdefs.IsOutOfRange(err):
		return New(ErrInfraOutOfRange).WithNTParam("Error", err.Error())
	case errdefs.IsNotImplemented(err):
		return New(ErrInfraNotImplemented).WithNTParam("Error", err.Error())
	case errdefs.IsInternal(err):
		return New(ErrInfraInternal).WithNTParam("Error", err.Error())
	case errdefs.IsUnavailable(err):
		return New(ErrInfraUnavailable).WithNTParam("Error", err.Error())
	case errdefs.IsDataLoss(err):
		return New(ErrInfraDataLoss).WithNTParam("Error", err.Error())
	case errdefs.IsUnauthorized(err):
		return New(ErrInfraUnauthorized).WithNTParam("Error", err.Error())
	}
	return New(err)
}
