package routes

import (
	echoSwagger "github.com/swaggo/echo-swagger"
	"go-authorization/constants"
	"go-authorization/docs"
	"go-authorization/lib"
)

// SwaggerRoutes
// @securityDefinitions.apikey Authorization
// @in header
// @name Authorization
// @schemes http https
// @tokenUrl http://localhost/api/login
// @basePath
// @contact.name Subrat
// @contact.email subrat.gyawali@gmail.com
type SwaggerRoutes struct {
	config  lib.Config
	logger  lib.Logger
	handler lib.HttpHandler
}

// NewSwaggerRoutes  creates new swagger routes
func NewSwaggerRoutes(
	config lib.Config,
	logger lib.Logger,
	handler lib.HttpHandler,
) SwaggerRoutes {
	return SwaggerRoutes{
		config:  config,
		logger:  logger,
		handler: handler,
	}
}

// Setup swagger routes
func (a SwaggerRoutes) Setup() {
	docs.SwaggerInfo.Title = a.config.Name
	docs.SwaggerInfo.Version = constants.Version
	docs.SwaggerInfo.BasePath = ""

	a.logger.Zap.Info("Setting up swagger routes")
	a.handler.Engine.GET("/swagger/*", echoSwagger.WrapHandler)
}
