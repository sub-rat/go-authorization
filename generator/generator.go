package generator

import (
	"fmt"
	"go-authorization/lib"
	"go-authorization/pkg/file"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

var controller = `
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

type ITemplateController interface {
	IController
}

type templateController struct {
	templateService services.TemplateService
	logger      lib.Logger
}

// NewTemplateController creates new template controller
func NewTemplateController(
	logger lib.Logger,
	templateService services.TemplateService,
) ITemplateController {
	return &templateController{
		logger:      logger,
		templateService: templateService,
	}
}

// Query
// @tags Template
// @summary Template Query
// @produce application/json
// @Security Authorization
// @param data query models.TemplateQueryParam true "TemplateQueryParam"
// @success 200 {object} echo_response.Response{data=models.TemplateQueryResult} "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/templates [get]
func (a *templateController) Query(ctx echo.Context) error {
	param := new(models.TemplateQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := a.templateService.Query(param)
	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}
	return echo_response.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// Get
// @tags Template
// @summary Template Get By ID
// @produce application/json
// @Security Authorization
// @param id path int true "template id"
// @success 200 {object} echo_response.Response{data=models.Template} "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/templates/{id} [get]
func (a *templateController) Get(ctx echo.Context) error {
	template, err := a.templateService.Get(ctx.Param("id"))
	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK, Data: template}.JSON(ctx)
}

// Create
// @tags Template
// @summary Template Create
// @produce application/json
// @Security Authorization
// @param data body models.Template true "Template"
// @success 200 {object} echo_response.Response "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/templates [post]
func (a *templateController) Create(ctx echo.Context) error {
	template := new(models.Template)
	if err := ctx.Bind(template); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	template.CreatedBy = claims.Username

	id, err := a.templateService.WithTrx(trxHandle).Create(template)
	if err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// Update
// @tags Template
// @summary Template Update By ID
// @produce application/json
// @Security Authorization
// @param id path int true "template id"
// @param data body models.Template true "Template"
// @success 200 {object} echo_response.Response "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/templates/{id} [put]
func (a *templateController) Update(ctx echo.Context) error {
	template := new(models.Template)
	if err := ctx.Bind(template); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.templateService.WithTrx(trxHandle).Update(ctx.Param("id"), template); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK}.JSON(ctx)
}

// Delete
// @tags Template
// @summary Template Delete By ID
// @produce application/json
// @Security Authorization
// @param id path int true "template id"
// @success 200 {object} echo_response.Response "ok"
// @failure 400 {object} echo_response.Response "bad request"
// @failure 500 {object} echo_response.Response "internal error"
// @router /api/v1/templates/{id} [delete]
func (a *templateController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.templateService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echo_response.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echo_response.Response{Code: http.StatusOK}.JSON(ctx)
}
`

func CreateFile(logger lib.Logger, name string) {
	repository = strings.ReplaceAll(repository, "Template", cases.Title(language.English).String(name))
	repository = strings.ReplaceAll(repository, "template", strings.ToLower(name))

	service = strings.ReplaceAll(service, "Template", cases.Title(language.English).String(name))
	service = strings.ReplaceAll(service, "template", strings.ToLower(name))

	controller = strings.ReplaceAll(controller, "Template", cases.Title(language.English).String(name))
	controller = strings.ReplaceAll(controller, "template", strings.ToLower(name))

	models = strings.ReplaceAll(models, "Template", cases.Title(language.English).String(name))
	models = strings.ReplaceAll(models, "template", strings.ToLower(name))
	models = strings.ReplaceAll(models, "'", "`")

	routes = strings.ReplaceAll(routes, "Template", cases.Title(language.English).String(name))
	routes = strings.ReplaceAll(routes, "template", strings.ToLower(name))

	WriteData(fmt.Sprintf("./models/%s.go", strings.ToLower(name)), models)
	WriteData(fmt.Sprintf("./internal/repository/%s_repository.go", strings.ToLower(name)), repository)
	WriteData(fmt.Sprintf("./internal/service/%s_service.go", strings.ToLower(name)), service)
	WriteData(fmt.Sprintf("./internal/controllers/%s_controller.go", strings.ToLower(name)), controller)
	WriteData(fmt.Sprintf("./internal/routes/%s_route.go", strings.ToLower(name)), routes)
}

func WriteData(fileName string, data string) {
	exists := file.IsExist(fileName)
	if exists {
		fmt.Printf("File already exists %s \n", fileName)
		return
	}
	f, err := file.Create(fileName)
	if err != nil {
		panic(err)
	}
	f.WriteString(data)
	f.Close()
}

var repository = `
package repository

import (
	"go-authorization/errors"
	"go-authorization/lib"
	"go-authorization/models"
	"gorm.io/gorm"
)

type ITemplateRepository interface {
	WithTrx(trxHandle *gorm.DB) ITemplateRepository
	Query(param *models.TemplateQueryParam) (*models.TemplateQueryResult, error)
	Get(id string) (*models.Template, error)
	Create(template *models.Template) error
	Update(id string, template *models.Template) error
	Delete(id string) error
}

// templateRepository database structure
type templateRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewTemplateRepository creates a new template repository
func NewTemplateRepository(db lib.Database, logger lib.Logger) ITemplateRepository {
	return templateRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a templateRepository) WithTrx(trxHandle *gorm.DB) ITemplateRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

// Query  gets all templates
func (a templateRepository) Query(param *models.TemplateQueryParam) (*models.TemplateQueryResult, error) {
	db := a.db.ORM.Model(&models.Template{})

	if v := param.QueryValue; v != "" {
		// Search query
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.Templates, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.TemplateQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a templateRepository) Get(id string) (*models.Template, error) {
	template := new(models.Template)

	if ok, err := QueryOne(a.db.ORM.Model(template).Where("id=?", id), template); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return template, nil
}

func (a templateRepository) Create(template *models.Template) error {
	result := a.db.ORM.Model(template).Create(template)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a templateRepository) Update(id string, template *models.Template) error {
	result := a.db.ORM.Model(template).Where("id=?", id).Updates(template)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a templateRepository) Delete(id string) error {
	template := new(models.Template)

	result := a.db.ORM.Model(template).Where("id=?", id).Delete(template)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
`

var service = `
package services

import (
	"go-authorization/internal/repository"
	"go-authorization/lib"
	"go-authorization/models"
	"go-authorization/pkg/uuid"
	"gorm.io/gorm"
)

// TemplateService service layer
type TemplateService struct {
	logger               lib.Logger
	templateRepository   repository.ITemplateRepository
}

// NewTemplateService creates a new templateservice
func NewTemplateService(
	logger lib.Logger,
	templateRepository repository.ITemplateRepository,
) TemplateService {
	return TemplateService{
		logger:               logger,
		templateRepository:   templateRepository,
	}
}

// WithTrx delegates transaction to repository database
func (a TemplateService) WithTrx(trxHandle *gorm.DB) TemplateService {
	a.templateRepository = a.templateRepository.WithTrx(trxHandle)
	return a
}

func (a TemplateService) Query(param *models.TemplateQueryParam) (templateQR *models.TemplateQueryResult, err error) {
	return a.templateRepository.Query(param)
}

func (a TemplateService) Get(id string) (*models.Template, error) {
	template, err := a.templateRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return template, nil
}


func (a TemplateService) Create(template *models.Template) (id string, err error) {
	if err = a.Check(template); err != nil {
		return
	}

	template.ID = uuid.MustString()
	if err = a.templateRepository.Create(template); err != nil {
		return
	}
	return template.ID, nil
}

func (a TemplateService) Update(id string, template *models.Template) error {
	oTemplate, err := a.Get(id)
	if err != nil {
		return err
	}

	template.ID = oTemplate.ID
	template.CreatedBy = oTemplate.CreatedBy
	template.CreatedAt = oTemplate.CreatedAt

	if err := a.templateRepository.Update(id, template); err != nil {
		return err
	}
	return nil
}

func (a TemplateService) Delete(id string) error {
	_, err := a.templateRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.templateRepository.Delete(id); err != nil {
		return err
	}
	return nil
}

func (a TemplateService) Check(item *models.Template) error {
	/* result, err := a.templateRepository.Query(&models.TemplateQueryParam{
	 	Add item to check before insertion
	 })

	if err != nil {
		return err
	} else if len(result.List) > 0 {
		return errors.TemplateAlreadyExists
	}
    */
	return nil
}
`

var models = `
package models

import (
	"go-authorization/models/database"
	"go-authorization/models/dto"
)

// Template
// Status - 1: Enable 0: Disable
type Template struct {
	database.Model
	ID        string    'gorm:"column:id;size:36;index;not null;" json:"id"'
	CreatedBy string    'gorm:"column:created_by;not null;" json:"created_by"'
}

type Templates []*Template

type TemplateQueryParam struct {
	dto.PaginationParam
	dto.OrderParam
	QueryValue    string  'query:"query_value"'
}

type TemplateQueryResult struct {
	List       Templates           'json:"list"'
	Pagination *dto.Pagination 'json:"pagination"'
}

func (a Templates) ToIDs() []string {
	ids := make([]string, len(a))
	for i, item := range a {
		ids[i] = item.ID
	}
	return ids
}
`

var routes = `
package routes

import (
	"go-authorization/internal/controllers"
	"go-authorization/lib"
)

type TemplateRoutes struct {
	logger         lib.Logger
	handler        lib.HttpHandler
	templateController controllers.ITemplateController
}

// NewTemplateRoutes creates new template routes
func NewTemplateRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	templateController controllers.ITemplateController,
) TemplateRoutes {
	return TemplateRoutes{
		handler:        handler,
		logger:         logger,
		templateController: templateController,
	}
}

// Setup template routes
func (a TemplateRoutes) Setup() {
	a.logger.Zap.Info("Setting up template routes")
	api := a.handler.RouterV1.Group("/templates")
	{
		api.GET("", a.templateController.Query)
		api.POST("", a.templateController.Create)
		api.GET("/:id", a.templateController.Get)
		api.PUT("/:id", a.templateController.Update)
		api.DELETE("/:id", a.templateController.Delete)
	}
}
`
