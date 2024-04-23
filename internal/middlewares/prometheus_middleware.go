package middlewares

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"go-authorization/lib"
)

// PrometheusMiddleware middleware for cors
type PrometheusMiddleware struct {
	handler lib.HttpHandler
	logger  lib.Logger
}

// NewPrometheusMiddleware creates new cors middleware
func NewPrometheusMiddleware(handler lib.HttpHandler, logger lib.Logger) PrometheusMiddleware {
	return PrometheusMiddleware{
		handler: handler,
		logger:  logger,
	}
}

func (a PrometheusMiddleware) Setup() {
	a.logger.Zap.Info("Setting up prometheus middleware")
	a.handler.Engine.Use(echoprometheus.NewMiddleware("go_authorization"))
}
