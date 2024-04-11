package repository

import (
	"go-authorization/errors"
	"go-authorization/lib"
	"go-authorization/models"
	"gorm.io/gorm"
)

type IMenuActionRepository interface {
	WithTrx(trxHandle *gorm.DB) IMenuActionRepository
	Query(param *models.MenuActionQueryParam) (*models.MenuActionQueryResult, error)
	Get(id string) (*models.MenuAction, error)
	Create(menuAction *models.MenuAction) error
	Update(id string, menuAction *models.MenuAction) error
	Delete(id string) error
	DeleteByMenuID(menuID string) error
}

// menuActionRepository database structure
type menuActionRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewMenuActionRepository creates a new menu action repository
func NewMenuActionRepository(db lib.Database, logger lib.Logger) IMenuActionRepository {
	return menuActionRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a menuActionRepository) WithTrx(trxHandle *gorm.DB) IMenuActionRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a menuActionRepository) Query(param *models.MenuActionQueryParam) (*models.MenuActionQueryResult, error) {
	db := a.db.ORM.Model(&models.MenuAction{})

	if v := param.MenuID; v != "" {
		db = db.Where("menu_id=?", v)
	}

	if v := param.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.MenuActions, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.MenuActionQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a menuActionRepository) Get(id string) (*models.MenuAction, error) {
	menuAction := new(models.MenuAction)

	if ok, err := QueryOne(a.db.ORM.Model(menuAction).Where("id=?", id), menuAction); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return menuAction, nil
}

func (a menuActionRepository) Create(menuAction *models.MenuAction) error {
	result := a.db.ORM.Model(menuAction).Create(menuAction)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a menuActionRepository) Update(id string, menuAction *models.MenuAction) error {
	result := a.db.ORM.Model(menuAction).Where("id=?", id).Updates(menuAction)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a menuActionRepository) Delete(id string) error {
	menuAction := new(models.MenuAction)

	result := a.db.ORM.Model(menuAction).Where("id=?", id).Delete(menuAction)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a menuActionRepository) DeleteByMenuID(menuID string) error {
	menuAction := new(models.MenuAction)

	result := a.db.ORM.Model(menuAction).Where("menu_id=?", menuID).Delete(menuAction)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
