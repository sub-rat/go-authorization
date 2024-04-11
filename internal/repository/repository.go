package repository

// Module exports dependency
type Module struct {
	MenuActionRepository         IMenuActionRepository
	MenuActionResourceRepository IMenuActionResourceRepository
	MenuRepository               IMenuRepository
	RoleMenuRepository           IRoleMenuRepository
	RoleRepository               IRoleRepository
	UserRepository               IUserRepository
	UserRoleRepository           IUserRoleRepository
}
