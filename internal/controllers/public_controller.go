package controllers

import (
	"go-authorization/constants"
	"go-authorization/errors"
	services "go-authorization/internal/service"
	"go-authorization/lib"
	"go-authorization/models/dto"
	"go-authorization/pkg/echo_response"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type IPublicController interface {
	SysRoutes(ctx echo.Context) error
	UserInfo(ctx echo.Context) error
	MenuTree(ctx echo.Context) error
	UserLogin(ctx echo.Context) error
	UserLogout(ctx echo.Context) error
}

type publicController struct {
	userService services.UserService
	authService services.AuthService
	logger      lib.Logger
}

// NewPublicController creates new public controller
func NewPublicController(
	userService services.UserService,
	authService services.AuthService,
	logger lib.Logger,
) IPublicController {
	return publicController{
		userService: userService,
		authService: authService,
		logger:      logger,
	}
}

type route struct {
	*echo.Route
	Name *struct{} `json:"name,omitempty"`
}

// SysRoutes
// @Tags Public
// @Summary SysRoutes
// @Produce application/json
// @Security Authorization
// @Success 200 {string} echo_response.Response "ok"
// @failure 400 {string} echo_response.Response "bad request"
// @failure 500 {string} echo_response.Response "internal error"
// @Router /api/v1/publics/sys/routes [get]
func (a publicController) SysRoutes(ctx echo.Context) error {
	routes := make([]*route, 0)
	for _, eRoute := range ctx.Echo().Routes() {
		// Only interfaces starting with /api/ are exposed
		if !strings.HasPrefix(eRoute.Path, "/api/") {
			continue
		}
		routes = append(routes, &route{Route: eRoute})
	}

	return echo_response.Response{Code: http.StatusOK, Data: routes}.JSON(ctx)
}

// UserInfo
// @Tags Public
// @Summary UserInfo
// @Security Authorization
// @Produce application/json
// @Success 200 {string} echo_response.Response{data=models.UserInfo} "ok"
// @failure 400 {string} echo_response.Response "bad request"
// @failure 500 {string} echo_response.Response "internal error"
// @Router /api/v1/publics/user [get]
func (a publicController) UserInfo(ctx echo.Context) error {
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)

	userinfo, err := a.userService.GetUserInfo(claims.ID)
	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK, Data: userinfo}.JSON(ctx)
}

// MenuTree
// @Tags Public
// @Summary UserMenuTree
// @Produce application/json
// @Security Authorization
// @Success 200 {string} echo_response.Response{data=models.MenuTrees} "ok"
// @failure 400 {string} echo_response.Response "bad request"
// @failure 500 {string} echo_response.Response "internal error"
// @Router /api/v1/publics/user/menutree [get]
func (a publicController) MenuTree(ctx echo.Context) error {
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)

	menuTrees, err := a.userService.GetUserMenuTrees(claims.ID)
	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK, Data: menuTrees}.JSON(ctx)
}

// UserLogin
// @Tags Public
// @Summary UserLogin
// @Produce application/json
// @Param data body dto.Login true "Login"
// @Success 200 {string} echo_response.Response "ok"
// @failure 400 {string} echo_response.Response "bad request"
// @failure 500 {string} echo_response.Response "internal error"
// @Router /api/v1/publics/user/login [post]
func (a publicController) UserLogin(ctx echo.Context) error {
	login := new(dto.Login)

	if err := ctx.Bind(login); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	//if !a.captcha.Verify(login.CaptchaID, login.CaptchaCode, false) {
	//	return echo_response.Response{Code: http.StatusBadRequest, Message: errors.CaptchaAnswerCodeNoMatch}.JSON(ctx)
	//}

	user, err := a.userService.Verify(login.Username, login.Password)
	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	token, err := a.authService.GenerateToken(user)
	if err != nil {
		return echo_response.Response{Code: http.StatusInternalServerError, Message: errors.AuthTokenGenerateFail}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK, Data: echo.Map{"token": token}}.JSON(ctx)
}

// UserLogout
// @Tags Public
// @Summary UserLogout
// @Produce application/json
// @Security Authorization
// @Success 200 {string} echo_response.Response "success"
// @Router /api/v1/publics/user/logout [post]
func (a publicController) UserLogout(ctx echo.Context) error {
	claims, ok := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	if ok {
		a.authService.DestroyToken(claims.Username)
	}

	return echo_response.Response{Code: http.StatusOK}.JSON(ctx)
}
