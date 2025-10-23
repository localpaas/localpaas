package s3storagehandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/s3storageuc"
)

type S3StorageHandler struct {
	*handler.BaseHandler
	authHandler *authhandler.AuthHandler
	s3StorageUC *s3storageuc.S3StorageUC
}

func NewS3StorageHandler(
	authHandler *authhandler.AuthHandler,
	s3StorageUC *s3storageuc.S3StorageUC,
) *S3StorageHandler {
	hdl := &S3StorageHandler{
		authHandler: authHandler,
		s3StorageUC: s3StorageUC,
	}
	return hdl
}
