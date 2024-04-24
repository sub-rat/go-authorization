package controllers

import (
	"go-authorization/constants"
	"go-authorization/errors"
	services "go-authorization/internal/service"
	"go-authorization/lib"
	"go-authorization/models"
	"go-authorization/models/dto"
	"go-authorization/pkg/echo_response"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

type IUserController interface {
	IController
	Enable(ctx echo.Context) error
	Disable(ctx echo.Context) error
}

type userController struct {
	userService services.UserService
	logger      lib.Logger
}

// NewUserController creates new user controller
func NewUserController(userService services.UserService, logger lib.Logger) IUserController {
	return &userController{
		userService: userService,
		logger:      logger,
	}
}

// Query
// @tags User
// @summary User Query
// @produce application/json
// @Security Authorization
// @param data query models.UserQueryParam true "UserQueryParam"
// @success 200 {object} echo_response.Response{data=models.UserQueryResult} "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/users [get]
func (a *userController) Query(ctx echo.Context) error {
	param := new(models.UserQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}
	if v := ctx.QueryParam("role_ids"); v != "" {
		param.RoleIDs = strings.Split(v, ",")
	}

	qr, err := a.userService.Query(param)
	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// Create
// @tags User
// @summary User Create
// @produce application/json
// @Security Authorization
// @param data body models.User true "User"
// @success 200 {object} echo_response.Response "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/users [post]
func (a *userController) Create(ctx echo.Context) error {
	user := new(models.User)
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)

	if err := ctx.Bind(user); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	} else if user.Password == "" {
		return echo_response.Response{Code: http.StatusBadRequest, Message: errors.UserPasswordRequired}.JSON(ctx)
	}

	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	user.CreatedBy = claims.Username

	qr, err := a.userService.WithTrx(trxHandle).Create(user)
	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// Get
// @tags User
// @summary User Get By ID
// @produce application/json
// @Security Authorization
// @param id path int true "user id"
// @success 200 {object} echo_response.Response{data=models.User} "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/users/{id} [get]
func (a *userController) Get(ctx echo.Context) error {
	user, err := a.userService.Get(ctx.Param("id"))
	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK, Data: user}.JSON(ctx)
}

// Update
// @tags User
// @summary User Update By ID
// @produce application/json
// @Security Authorization
// @param id path int true "user id"
// @param data body models.User true "User"
// @success 200 {object} echo_response.Response "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/users/{id} [put]
func (a *userController) Update(ctx echo.Context) error {
	user := new(models.User)
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)

	if err := ctx.Bind(user); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	err := a.userService.WithTrx(trxHandle).Update(ctx.Param("id"), user)
	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK}.JSON(ctx)
}

// Delete
// @tags User
// @summary User Delete By ID
// @produce application/json
// @Security Authorization
// @param id path int true "user id"
// @success 200 {object} echo_response.Response "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/users/{id} [delete]
func (a *userController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	err := a.userService.WithTrx(trxHandle).Delete(ctx.Param("id"))
	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK}.JSON(ctx)
}

// Enable
// @tags User
// @summary User Enable By ID
// @produce application/json
// @Security Authorization
// @param id path int true "user id"
// @success 200 {object} echo_response.Response "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/users/{id}/enable [patch]
func (a *userController) Enable(ctx echo.Context) error {
	err := a.userService.UpdateStatus(ctx.Param("id"), 1)
	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK}.JSON(ctx)
}

// Disable
// @tags User
// @summary User Disable By ID
// @produce application/json
// @Security Authorization
// @param id path int true "user id"
// @success 200 {object} echo_response.Response "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/users/{id}/disable [patch]
func (a *userController) Disable(ctx echo.Context) error {
	err := a.userService.UpdateStatus(ctx.Param("id"), -1)
	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK}.JSON(ctx)
}
