package repository

import (
	"go-authorization/errors"
	"go-authorization/lib"
	"go-authorization/models"
	"gorm.io/gorm"
)

type IMenuActionResourceRepository interface {
	WithTrx(trxHandle *gorm.DB) IMenuActionResourceRepository
	Query(param *models.MenuActionResourceQueryParam) (*models.MenuActionResourceQueryResult, error)
	Get(id string) (*models.MenuActionResource, error)
	Create(menuActionResource *models.MenuActionResource) error
	Update(id string, menuActionResource *models.MenuActionResource) error
	Delete(id string) error
	DeleteByActionID(actionID string) error
	DeleteByMenuID(menuID string) error
}

// menuActionRepository database structure
type menuActionResourceRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewMenuActionResourceRepository creates a new menu action resource repository
func NewMenuActionResourceRepository(db lib.Database, logger lib.Logger) IMenuActionResourceRepository {
	return menuActionResourceRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a menuActionResourceRepository) WithTrx(trxHandle *gorm.DB) IMenuActionResourceRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a menuActionResourceRepository) Query(param *models.MenuActionResourceQueryParam) (*models.MenuActionResourceQueryResult, error) {
	db := a.db.ORM.Model(&models.MenuActionResource{})

	if v := param.MenuID; v != "" {
		subQuery := a.db.ORM.Model(&models.MenuAction{}).
			Where("menu_id=?", v).
			Select("id")

		db = db.Where("action_id IN (?)", subQuery)
	}

	if v := param.MenuIDs; len(v) > 0 {
		subQuery := a.db.ORM.Model(&models.MenuAction{}).
			Where("menu_id IN (?)", v).
			Select("id")

		db = db.Where("action_id IN (?)", subQuery)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.MenuActionResources, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.MenuActionResourceQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a menuActionResourceRepository) Get(id string) (*models.MenuActionResource, error) {
	menuActionResource := new(models.MenuActionResource)

	if ok, err := QueryOne(a.db.ORM.Model(menuActionResource).Where("id=?", id), menuActionResource); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return menuActionResource, nil
}

func (a menuActionResourceRepository) Create(menuActionResource *models.MenuActionResource) error {
	result := a.db.ORM.Model(menuActionResource).Create(menuActionResource)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a menuActionResourceRepository) Update(id string, menuActionResource *models.MenuActionResource) error {
	result := a.db.ORM.Model(menuActionResource).Where("id=?", id).Updates(menuActionResource)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a menuActionResourceRepository) Delete(id string) error {
	menuActionResource := new(models.MenuActionResource)

	result := a.db.ORM.Model(menuActionResource).Where("id=?", id).Delete(menuActionResource)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a menuActionResourceRepository) DeleteByActionID(actionID string) error {
	menuActionResource := new(models.MenuActionResource)

	result := a.db.ORM.Model(menuActionResource).Where("action_id=?", actionID).Delete(menuActionResource)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a menuActionResourceRepository) DeleteByMenuID(menuID string) error {
	menuAction := new(models.MenuAction)
	menuActionResource := new(models.MenuActionResource)

	subQuery := a.db.ORM.Model(menuAction).
		Where("menu_id=?", menuID).Select("id")

	result := a.db.ORM.Model(menuActionResource).
		Where("action_id IN (?)", subQuery).Delete(menuActionResource)

	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
