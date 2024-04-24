package controllers

import (
	"go-authorization/constants"
	services "go-authorization/internal/service"
	"go-authorization/lib"
	"go-authorization/models"
	"go-authorization/models/dto"
	"go-authorization/pkg/echo_response"
	"net/http"

	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

type IMenuController interface {
	IController
	Enable(ctx echo.Context) error
	Disable(ctx echo.Context) error
	GetActions(ctx echo.Context) error
	UpdateActions(ctx echo.Context) error
}

type menuController struct {
	menuService services.MenuService
	logger      lib.Logger
}

// NewMenuController creates new menu controller
func NewMenuController(
	logger lib.Logger,
	menuService services.MenuService,
) IMenuController {
	return &menuController{
		logger:      logger,
		menuService: menuService,
	}
}

// Query
// @tags Menu
// @summary Menu Query
// @produce application/json
// @Security Authorization
// @param data query models.MenuQueryParam true "MenuQueryParam"
// @success 200 {object} echo_response.Response{data=models.MenuQueryResult} "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/menus [get]
func (a *menuController) Query(ctx echo.Context) error {
	param := new(models.MenuQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := a.menuService.Query(param)
	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	if param.Tree {
		return echo_response.Response{Code: http.StatusOK, Data: qr.List.ToMenuTrees()}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// Get
// @tags Menu
// @summary Menu Get By ID
// @produce application/json
// @Security Authorization
// @param id path int true "menu id"
// @success 200 {object} echo_response.Response{data=models.Menu} "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/menus/{id} [get]
func (a *menuController) Get(ctx echo.Context) error {
	menu, err := a.menuService.Get(ctx.Param("id"))
	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK, Data: menu}.JSON(ctx)
}

// Create
// @tags Menu
// @summary Menu Create
// @produce application/json
// @Security Authorization
// @param data body models.Menu true "Menu"
// @success 200 {object} echo_response.Response "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/menus [post]
func (a *menuController) Create(ctx echo.Context) error {
	menu := new(models.Menu)
	if err := ctx.Bind(menu); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	menu.CreatedBy = claims.Username

	id, err := a.menuService.WithTrx(trxHandle).Create(menu)
	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// Update
// @tags Menu
// @summary Menu Update By ID
// @produce application/json
// @Security Authorization
// @param id path int true "menu id"
// @param data body models.Menu true "Menu"
// @success 200 {object} echo_response.Response "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/menus/{id} [put]
func (a *menuController) Update(ctx echo.Context) error {
	menu := new(models.Menu)
	if err := ctx.Bind(menu); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.menuService.WithTrx(trxHandle).Update(ctx.Param("id"), menu); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK}.JSON(ctx)
}

// Delete
// @tags Menu
// @summary Menu Delete By ID
// @produce application/json
// @Security Authorization
// @param id path int true "menu id"
// @success 200 {object} echo_response.Response "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/menus/{id} [delete]
func (a *menuController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.menuService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK}.JSON(ctx)
}

// Enable
// @tags Menu
// @summary Menu Enable By ID
// @produce application/json
// @Security Authorization
// @param id path int true "menu id"
// @success 200 {object} echo_response.Response "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/menus/{id}/enable [patch]
func (a *menuController) Enable(ctx echo.Context) error {
	if err := a.menuService.UpdateStatus(ctx.Param("id"), 1); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK}.JSON(ctx)
}

// Disable
// @tags Menu
// @summary Menu Disable By ID
// @produce application/json
// @Security Authorization
// @param id path int true "menu id"
// @success 200 {object} echo_response.Response "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/menus/{id}/disable [patch]
func (a *menuController) Disable(ctx echo.Context) error {
	if err := a.menuService.UpdateStatus(ctx.Param("id"), -1); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK}.JSON(ctx)
}

// GetActions
// @tags Menu
// @summary MenuActions Get By menuID
// @produce application/json
// @Security Authorization
// @param id path int true "menu id"
// @success 200 {object} echo_response.Response{data=models.MenuActions} "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/menus/{id}/actions [get]
func (a *menuController) GetActions(ctx echo.Context) error {
	actions, err := a.menuService.GetMenuActions(ctx.Param("id"))
	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK, Data: actions}.JSON(ctx)
}

// UpdateActions
// @tags Menu
// @summary Menu Actions Update By menuID
// @produce application/json
// @Security Authorization
// @param id path int true "menu id"
// @param data body models.MenuActions true "Menu"
// @success 200 {object} echo_response.Response "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/menus/{id}/actions [put]
func (a *menuController) UpdateActions(ctx echo.Context) error {
	actions := make(models.MenuActions, 0)
	if err := ctx.Bind(&actions); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.menuService.WithTrx(trxHandle).UpdateActions(ctx.Param("id"), actions); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK}.JSON(ctx)
}
