package bootstrap

import "go-authorization/internal/repository"

func (app *Application) ProvideRepository() {
	app.Repository = repository.Module{
		MenuActionRepository:         repository.NewMenuActionRepository(app.Database, app.Logger),
		MenuActionResourceRepository: repository.NewMenuActionResourceRepository(app.Database, app.Logger),
		MenuRepository:               repository.NewMenuRepository(app.Database, app.Logger),
		RoleMenuRepository:           repository.NewRoleMenuRepository(app.Database, app.Logger),
		RoleRepository:               repository.NewRoleRepository(app.Database, app.Logger),
		UserRepository:               repository.NewUserRepository(app.Database, app.Logger),
		UserRoleRepository:           repository.NewUserRoleRepository(app.Database, app.Logger),
	}
}
