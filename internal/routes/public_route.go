package routes

import (
	"go-authorization/internal/controllers"
	"go-authorization/lib"
)

type PublicRoutes struct {
	logger           lib.Logger
	handler          lib.HttpHandler
	publicController controllers.IPublicController
}

// NewPublicRoutes creates new public routes
func NewPublicRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	publicController controllers.IPublicController,
) PublicRoutes {
	return PublicRoutes{
		handler:          handler,
		logger:           logger,
		publicController: publicController,
	}
}

// Setup public routes
func (a PublicRoutes) Setup() {
	a.logger.Zap.Info("Setting up public routes")
	api := a.handler.RouterV1.Group("/publics")
	{
		api.GET("/user", a.publicController.UserInfo)
		api.POST("/user/login", a.publicController.UserLogin)
		api.POST("/user/logout", a.publicController.UserLogout)
		api.GET("/user/menutree", a.publicController.MenuTree)

		// sys routes
		api.GET("/sys/routes", a.publicController.SysRoutes)

	}
}
