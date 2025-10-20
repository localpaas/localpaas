package recovery

import (
	"io"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
)

// Recovery create a middleware for recovering from panic
func Recovery(cfg *config.Config) gin.HandlerFunc {
	// In production, use `nil` as writer to prevent Gin to log sensitive information
	// of the request to the default stderr.
	var writer io.Writer
	if !cfg.IsProdEnv() {
		writer = gin.DefaultErrorWriter
	}

	return gin.CustomRecoveryWithWriter(writer, func(ctx *gin.Context, recover any) {
		err := apperrors.New(apperrors.ErrInternalServer).
			WithMsgLog("recovered from panic: %v", recover)
		(&handler.BaseHandler{}).RenderError(ctx, err)
		ctx.Abort()
	})
}
