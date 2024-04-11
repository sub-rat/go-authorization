package bootstrap

import (
	"go-authorization/internal/controllers"
	"go-authorization/internal/routes"
)

func (app *Application) ProvideControllers() {
	app.Controllers = controllers.Module{
		MenuController:   controllers.NewMenuController(app.Logger, app.Services.MenuService),
		PublicController: controllers.NewPublicController(app.Services.UserService, app.Services.AuthService, app.Logger),
		RoleController:   controllers.NewRoleController(app.Logger, app.Services.RoleService),
		UserController:   controllers.NewUserController(app.Services.UserService, app.Logger),
	}
}

func (app *Application) ProvideRoutes() {
	app.Routes = routes.NewRoutes(
		routes.NewMenuRoutes(app.Logger, app.Handler, app.Controllers.MenuController),
		routes.NewPprofRoutes(app.Logger, app.Handler),
		routes.NewSwaggerRoutes(app.Config, app.Logger, app.Handler),
		routes.NewPublicRoutes(app.Logger, app.Handler, app.Controllers.PublicController),
		routes.NewUserRoutes(app.Logger, app.Handler, app.Controllers.UserController),
		routes.NewRoleRoutes(app.Logger, app.Handler, app.Controllers.RoleController),
	)
}
