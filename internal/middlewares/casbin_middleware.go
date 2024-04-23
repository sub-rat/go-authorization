package middlewares

import (
	"github.com/labstack/echo/v4"
	"go-authorization/constants"
	services "go-authorization/internal/service"
	"go-authorization/lib"
	"go-authorization/models/dto"
	"go-authorization/pkg/echo_response"
	"net/http"
)

// CasbinMiddleware PrometheusMiddleware middleware for cors
type CasbinMiddleware struct {
	handler       lib.HttpHandler
	logger        lib.Logger
	config        lib.Config
	casbinService services.CasbinService
}

// NewCasbinMiddleware NewCorsMiddleware creates new cors middleware
func NewCasbinMiddleware(
	handler lib.HttpHandler,
	logger lib.Logger,
	config lib.Config,
	casbinService services.CasbinService,
) CasbinMiddleware {
	return CasbinMiddleware{
		handler:       handler,
		logger:        logger,
		config:        config,
		casbinService: casbinService,
	}
}

func (a CasbinMiddleware) core() echo.MiddlewareFunc {
	prefixes := a.config.Casbin.IgnorePathPrefixes

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			request := ctx.Request()
			if isIgnorePath(request.URL.Path, prefixes...) {
				return next(ctx)
			}

			p := ctx.Request().URL.Path
			m := ctx.Request().Method
			claims, ok := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
			if !ok {
				return echo_response.Response{Code: http.StatusUnauthorized}.JSON(ctx)
			}
			a.logger.Zap.Info(p, m, claims)
			if ok, err := a.casbinService.Enforcer.Enforce(claims.ID, p, m); err != nil {
				return echo_response.Response{Code: http.StatusForbidden, Message: err}.JSON(ctx)
			} else if !ok {
				return echo_response.Response{Code: http.StatusForbidden}.JSON(ctx)
			}

			return next(ctx)
		}
	}
}

func (a CasbinMiddleware) Setup() {
	if !a.config.Casbin.Enable {
		return
	}

	a.logger.Zap.Info("Setting up casbin middleware")
	a.handler.Engine.Use(a.core())
}
