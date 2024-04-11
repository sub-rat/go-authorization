package controllers

import (
	"github.com/labstack/echo/v4"
)

// Module exported for initializing application
type Module struct {
	MenuController   IMenuController
	PublicController IPublicController
	RoleController   IRoleController
	UserController   IUserController
}

type IController interface {
	Query(ctx echo.Context) error
	Get(ctx echo.Context) error
	Create(ctx echo.Context) error
	Update(ctx echo.Context) error
	Delete(ctx echo.Context) error
}
