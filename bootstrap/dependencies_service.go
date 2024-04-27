package bootstrap

import services "go-authorization/internal/service"

func (app *Application) ProvideService() {
	// Auth Service
	app.Services = services.Module{
		AuthService: services.NewAuthService(app.Redis, app.Config),
	}

	// Casbin Service
	app.Services.CasbinService = services.NewCasbinService(app.Logger, app.Config, app.Repository.UserRepository,
		app.Repository.UserRoleRepository, app.Repository.RoleRepository, app.Repository.RoleMenuRepository,
		app.Repository.MenuActionResourceRepository)

	// Menu Service
	app.Services.MenuService = services.NewMenuService(app.Logger, app.Repository.MenuRepository,
		app.Repository.MenuActionRepository, app.Repository.MenuActionResourceRepository)

	// Role Service
	app.Services.RoleService = services.NewRoleService(app.Logger, app.Services.CasbinService,
		app.Repository.UserRepository, app.Repository.RoleRepository, app.Repository.RoleMenuRepository,
		app.Repository.MenuRepository, app.Repository.MenuActionRepository)

	// User Service
	app.Services.UserService = services.NewUserService(app.Logger,
		app.Repository.UserRepository, app.Repository.UserRoleRepository, app.Repository.RoleRepository,
		app.Repository.RoleMenuRepository, app.Repository.MenuRepository, app.Repository.MenuActionRepository,
		app.Services.CasbinService, app.Config)

}
