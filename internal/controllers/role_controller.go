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

type IRoleController interface {
	IController
	GetAll(ctx echo.Context) error
	Enable(ctx echo.Context) error
	Disable(ctx echo.Context) error
}

type roleController struct {
	logger      lib.Logger
	roleService services.RoleService
}

// NewRoleController creates new role controller
func NewRoleController(
	logger lib.Logger,
	roleService services.RoleService,
) IRoleController {
	return &roleController{
		logger:      logger,
		roleService: roleService,
	}
}

// Query
// @tags Role
// @summary Role Query
// @produce application/json
// @Security Authorization
// @param data query models.RoleQueryParam true "RoleQueryParam"
// @success 200 {object} echo_response.Response{data=models.RoleQueryResult} "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/roles [get]
func (a *roleController) Query(ctx echo.Context) error {
	param := new(models.RoleQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := a.roleService.Query(param)
	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// GetAll
// @tags Role
// @summary Role Get All
// @produce application/json
// @Security Authorization
// @param data query models.RoleQueryParam true "RoleQueryParam"
// @success 200 {object} echo_response.Response{data=models.Roles} "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/roles [get]
func (a *roleController) GetAll(ctx echo.Context) error {
	qr, err := a.roleService.Query(&models.RoleQueryParam{
		PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})

	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK, Data: qr.List}.JSON(ctx)
}

// Get
// @tags Role
// @summary Role Get By ID
// @produce application/json
// @Security Authorization
// @param id path int true "role id"
// @success 200 {object} echo_response.Response{data=models.Role} "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/roles/{id} [get]
func (a *roleController) Get(ctx echo.Context) error {
	role, err := a.roleService.Get(ctx.Param("id"))
	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK, Data: role}.JSON(ctx)
}

// Create
// @tags Role
// @summary Role Create
// @produce application/json
// @Security Authorization
// @param data body models.Role true "Role"
// @success 200 {object} echo_response.Response "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/roles [post]
func (a *roleController) Create(ctx echo.Context) error {
	role := new(models.Role)
	if err := ctx.Bind(role); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	role.CreatedBy = claims.Username

	id, err := a.roleService.WithTrx(trxHandle).Create(role)
	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// Update
// @tags Role
// @summary Role Update By ID
// @produce application/json
// @Security Authorization
// @param id path int true "role id"
// @param data body models.Role true "Role"
// @success 200 {object} echo_response.Response "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/roles/{id} [put]
func (a *roleController) Update(ctx echo.Context) error {
	role := new(models.Role)
	if err := ctx.Bind(role); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.roleService.WithTrx(trxHandle).Update(ctx.Param("id"), role); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK}.JSON(ctx)
}

// Delete
// @tags Role
// @summary Role Delete By ID
// @produce application/json
// @Security Authorization
// @param id path int true "role id"
// @success 200 {object} echo_response.Response "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/roles/{id} [delete]
func (a *roleController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.roleService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK}.JSON(ctx)
}

// Enable
// @tags Role
// @summary Role Enable By ID
// @produce application/json
// @Security Authorization
// @param id path int true "role id"
// @success 200 {object} echo_response.Response "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/roles/{id}/enable [patch]
func (a *roleController) Enable(ctx echo.Context) error {
	if err := a.roleService.UpdateStatus(ctx.Param("id"), 1); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK}.JSON(ctx)
}

// Disable
// @tags Role
// @summary Role Disable By ID
// @produce application/json
// @Security Authorization
// @param id path int true "role id"
// @success 200 {object} echo_response.Response "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/roles/{id}/disable [patch]
func (a *roleController) Disable(ctx echo.Context) error {
	if err := a.roleService.UpdateStatus(ctx.Param("id"), -1); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK}.JSON(ctx)
}
