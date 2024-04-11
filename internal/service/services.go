package services

type Module struct {
	AuthService   AuthService
	CasbinService CasbinService
	MenuService   MenuService
	RoleService   RoleService
	UserService   UserService
}
