package server

import (
	"context"
	"net/http"
	"time"

	ginlogger "github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/middleware/cors"
	loggermiddleware "github.com/localpaas/localpaas/localpaas_app/interface/api/middleware/logger"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/middleware/recovery"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/middleware/secureheaders"
)

type Server interface {
	Start() error
	Stop(context.Context) error
	GetAddress() string
}

type HTTPServer struct {
	*http.Server
	config          *config.Config
	engine          *gin.Engine
	handlerRegistry *HandlerRegistry
	logger          logging.Logger
}

// NewHTTPServer Create new LocalPaas app
// @title LocalPaas Backend
// @version 0.1
// @description LocalPaas Backend
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /v1
// @securityDefinitions.basic BasicAuth
func NewHTTPServer(
	config *config.Config,
	logger logging.Logger,
	handlerRegistry *HandlerRegistry,
) Server {
	engine := gin.New()
	s := &HTTPServer{
		config:          config,
		engine:          engine,
		handlerRegistry: handlerRegistry,
		logger:          logger,
		Server: &http.Server{
			Addr:           config.HTTPServer.BindingAddress(),
			ReadTimeout:    10 * time.Second, //nolint:mnd
			WriteTimeout:   10 * time.Second, //nolint:mnd
			MaxHeaderBytes: 1 << 20,          //nolint:mnd
			Handler:        engine,
		},
	}

	// Configures middlewares
	engine.Use(
		recovery.Recovery(config),
		requestid.New(),
		loggermiddleware.Logger(logger),
		secureheaders.SecureHeaders,
		cors.CORS(config),
	)

	if !config.IsProdEnv() {
		engine.Use(ginlogger.SetLogger())
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	s.registerRoutes()

	return s
}

func (s *HTTPServer) Start() error {
	err := s.ListenAndServe()
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	err := s.Shutdown(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (s *HTTPServer) GetAddress() string {
	return s.Addr
}
