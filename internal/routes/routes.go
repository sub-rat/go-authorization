package routes

import (
	"go-authorization/lib"
)

// Module exports dependency to container
type Module struct {
	MenuRoutes    MenuRoutes
	PprofRoutes   PprofRoutes
	PublicRoutes  PublicRoutes
	RoleRoutes    RoleRoutes
	SwaggerRoutes SwaggerRoutes
	UserRoutes    UserRoutes
}

// Route interface
type Route interface {
	Setup()
}

// NewRoutes sets up routes
func NewRoutes(
	menuRoutes MenuRoutes,
	pprofRoutes PprofRoutes,
	swaggerRoutes SwaggerRoutes,
	publicRoutes PublicRoutes,
	userRoutes UserRoutes,
	roleRoutes RoleRoutes,
) Module {
	return Module{
		MenuRoutes:    menuRoutes,
		PprofRoutes:   pprofRoutes,
		SwaggerRoutes: swaggerRoutes,
		PublicRoutes:  publicRoutes,
		UserRoutes:    userRoutes,
		RoleRoutes:    roleRoutes,
	}
}

// Setup all the route
func (a Module) Setup(logger lib.Logger) {
	logger.DesugarZap.Info("Setting Up Routes")
	routes := []Route{
		a.MenuRoutes,
		a.PprofRoutes,
		a.PublicRoutes,
		a.RoleRoutes,
		a.SwaggerRoutes,
		a.UserRoutes,
	}

	for _, route := range routes {
		route.Setup()
	}
}
