package bootstrap

import services "go-authorization/internal/service"

func (app *Application) ProvideService() {
	app.Services = services.Module{
		AuthService: services.NewAuthService(app.Redis, app.Config),
		CasbinService: services.NewCasbinService(app.Logger, app.Config, app.Repository.UserRepository,
			app.Repository.UserRoleRepository, app.Repository.RoleRepository, app.Repository.RoleMenuRepository,
			app.Repository.MenuActionResourceRepository),
		MenuService: services.NewMenuService(app.Logger, app.Repository.MenuRepository,
			app.Repository.MenuActionRepository, app.Repository.MenuActionResourceRepository),
		RoleService: services.NewRoleService(app.Logger, app.Services.CasbinService,
			app.Repository.UserRepository, app.Repository.RoleRepository, app.Repository.RoleMenuRepository,
			app.Repository.MenuRepository, app.Repository.MenuActionRepository),
		UserService: services.NewUserService(app.Logger,
			app.Repository.UserRepository, app.Repository.UserRoleRepository, app.Repository.RoleRepository,
			app.Repository.RoleMenuRepository, app.Repository.MenuRepository, app.Repository.MenuActionRepository,
			app.Services.CasbinService, app.Config),
	}
}
