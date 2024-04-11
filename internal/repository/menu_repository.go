package repository

import (
	"go-authorization/errors"
	"go-authorization/lib"
	"go-authorization/models"
	"gorm.io/gorm"
)

type IMenuRepository interface {
	WithTrx(trxHandle *gorm.DB) IMenuRepository
	Query(param *models.MenuQueryParam) (*models.MenuQueryResult, error)
	Get(id string) (*models.Menu, error)
	Create(menu *models.Menu) error
	Update(id string, menu *models.Menu) error
	Delete(id string) error
	UpdateStatus(id string, status int) error
	UpdateParentPath(id string, parentPath string) error
}

// menuRepository database structure
type menuRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewMenuRepository creates a new menu repository
func NewMenuRepository(db lib.Database, logger lib.Logger) IMenuRepository {
	return menuRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a menuRepository) WithTrx(trxHandle *gorm.DB) IMenuRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a menuRepository) Query(param *models.MenuQueryParam) (*models.MenuQueryResult, error) {
	db := a.db.ORM.Model(&models.Menu{})

	if v := param.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	if v := param.Name; v != "" {
		db = db.Where("name=?", v)
	}

	if v := param.ParentID; v != "" {
		db = db.Where("parent_id=?", v)
	}

	if v := param.PrefixParentPath; v != "" {
		db = db.Where("parent_path LIKE ?", v+"%")
	}

	if v := param.Hidden; v != 0 {
		db = db.Where("show_status=?", v)
	}

	if v := param.Status; v != 0 {
		db = db.Where("status=?", v)
	}

	if v := param.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("name LIKE ? OR remark LIKE ?", v, v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.Menus, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.MenuQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a menuRepository) Get(id string) (*models.Menu, error) {
	menu := new(models.Menu)

	if ok, err := QueryOne(a.db.ORM.Model(menu).Where("id=?", id), menu); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return menu, nil
}

func (a menuRepository) Create(menu *models.Menu) error {
	result := a.db.ORM.Model(menu).Create(menu)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a menuRepository) Update(id string, menu *models.Menu) error {
	result := a.db.ORM.Model(menu).Where("id=?", id).Updates(menu)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a menuRepository) Delete(id string) error {
	menu := new(models.Menu)

	result := a.db.ORM.Model(menu).Where("id=?", id).Delete(menu)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a menuRepository) UpdateStatus(id string, status int) error {
	menu := new(models.Menu)

	result := a.db.ORM.Model(menu).Where("id=?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a menuRepository) UpdateParentPath(id string, parentPath string) error {
	menu := new(models.Menu)

	result := a.db.ORM.Model(menu).Where("id=?", id).Update("parent_path", parentPath)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
