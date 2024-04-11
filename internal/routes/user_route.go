package routes

import (
	"go-authorization/internal/controllers"
	"go-authorization/lib"
)

type UserRoutes struct {
	logger         lib.Logger
	handler        lib.HttpHandler
	userController controllers.IUserController
}

// NewUserRoutes creates new user routes
func NewUserRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	userController controllers.IUserController,
) UserRoutes {
	return UserRoutes{
		handler:        handler,
		logger:         logger,
		userController: userController,
	}
}

// Setup user routes
func (a UserRoutes) Setup() {
	a.logger.Zap.Info("Setting up user routes")
	api := a.handler.RouterV1.Group("/users")
	{
		api.GET("", a.userController.Query)
		api.POST("", a.userController.Create)
		api.GET("/:id", a.userController.Get)
		api.PUT("/:id", a.userController.Update)
		api.DELETE("/:id", a.userController.Delete)
		api.POST("/:id/enable", a.userController.Enable)
		api.POST("/:id/disable", a.userController.Disable)
	}
}
