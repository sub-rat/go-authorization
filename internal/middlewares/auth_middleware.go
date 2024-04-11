package middlewares

import (
	"go-authorization/constants"
	services "go-authorization/internal/service"
	"go-authorization/lib"
	"go-authorization/pkg/echo_response"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware middleware for cors
type AuthMiddleware struct {
	config      lib.Config
	handler     lib.HttpHandler
	logger      lib.Logger
	authService services.AuthService
}

// NewCorsMiddleware creates new cors middleware
func NewAuthMiddleware(
	handler lib.HttpHandler,
	logger lib.Logger,
	config lib.Config,
	authService services.AuthService,
) AuthMiddleware {
	return AuthMiddleware{
		config:      config,
		handler:     handler,
		logger:      logger,
		authService: authService,
	}
}

func (a AuthMiddleware) core() echo.MiddlewareFunc {
	prefixes := a.config.Auth.IgnorePathPrefixes

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			request := ctx.Request()
			if isIgnorePath(request.URL.Path, prefixes...) {
				return next(ctx)
			}

			var (
				auth   = request.Header.Get("Authorization")
				prefix = "Bearer "
				token  string
			)

			if auth != "" && strings.HasPrefix(auth, prefix) {
				token = auth[len(prefix):]
			}

			claims, err := a.authService.ParseToken(token)
			if err != nil {
				return echo_response.Response{Code: http.StatusUnauthorized, Message: err}.JSON(ctx)
			}

			ctx.Set(constants.CurrentUser, claims)
			return next(ctx)
		}
	}
}

func (a AuthMiddleware) Setup() {
	if !a.config.Auth.Enable {
		return
	}

	a.logger.Zap.Info("Setting up auth middleware")
	a.handler.Engine.Use(a.core())
}
