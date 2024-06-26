package routes

import (
	"go-authorization/internal/controllers"
	"go-authorization/lib"
)

type RoleRoutes struct {
	logger         lib.Logger
	handler        lib.HttpHandler
	roleController controllers.IRoleController
}

// NewRoleRoutes creates new role routes
func NewRoleRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	roleController controllers.IRoleController,
) RoleRoutes {
	return RoleRoutes{
		handler:        handler,
		logger:         logger,
		roleController: roleController,
	}
}

// Setup role routes
func (a RoleRoutes) Setup() {
	a.logger.Zap.Info("Setting up role routes")
	api := a.handler.RouterV1.Group("/roles")
	{
		api.GET("", a.roleController.Query)
		api.GET(".all", a.roleController.GetAll)

		api.POST("", a.roleController.Create)
		api.GET("/:id", a.roleController.Get)
		api.PUT("/:id", a.roleController.Update)
		api.DELETE("/:id", a.roleController.Delete)
		api.PATCH("/:id/enable", a.roleController.Enable)
		api.PATCH("/:id/disable", a.roleController.Disable)
	}
}
