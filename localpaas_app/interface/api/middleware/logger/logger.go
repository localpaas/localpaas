package logger

import (
	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
)

func Logger(logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(logging.LoggerCtxKey, logger)

		c.Next()
	}
}
