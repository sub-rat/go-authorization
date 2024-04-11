package repository

import (
	"go-authorization/errors"
	"go-authorization/lib"
	"go-authorization/models"
	"gorm.io/gorm"
)

type IRoleMenuRepository interface {
	WithTrx(trxHandle *gorm.DB) IRoleMenuRepository
	Query(param *models.RoleMenuQueryParam) (*models.RoleMenuQueryResult, error)
	Get(id string) (*models.RoleMenu, error)
	Create(roleMenu *models.RoleMenu) error
	Update(id string, roleMenu *models.RoleMenu) error
	Delete(id string) error
	DeleteByRoleID(id string) error
}

// roleMenuRepository database structure
type roleMenuRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewRoleMenuRepository creates a new role menu repository
func NewRoleMenuRepository(db lib.Database, logger lib.Logger) IRoleMenuRepository {
	return roleMenuRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a roleMenuRepository) WithTrx(trxHandle *gorm.DB) IRoleMenuRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a roleMenuRepository) Query(param *models.RoleMenuQueryParam) (*models.RoleMenuQueryResult, error) {
	db := a.db.ORM.Model(&models.RoleMenu{})

	if v := param.RoleID; v != "" {
		db = db.Where("role_id=?", v)
	}

	if v := param.RoleIDs; len(v) > 0 {
		db = db.Where("role_id IN (?)", v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make([]*models.RoleMenu, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.RoleMenuQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a roleMenuRepository) Get(id string) (*models.RoleMenu, error) {
	roleMenu := new(models.RoleMenu)

	if ok, err := QueryOne(a.db.ORM.Model(roleMenu).Where("id=?", id), roleMenu); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return roleMenu, nil
}

func (a roleMenuRepository) Create(roleMenu *models.RoleMenu) error {
	result := a.db.ORM.Model(roleMenu).Create(roleMenu)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a roleMenuRepository) Update(id string, roleMenu *models.RoleMenu) error {
	result := a.db.ORM.Model(roleMenu).Where("id=?", id).Updates(roleMenu)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a roleMenuRepository) Delete(id string) error {
	roleMenu := new(models.RoleMenu)

	result := a.db.ORM.Model(roleMenu).Where("id=?", id).Delete(roleMenu)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a roleMenuRepository) DeleteByRoleID(id string) error {
	roleMenu := new(models.RoleMenu)

	result := a.db.ORM.Model(roleMenu).Where("role_id=?", id).Delete(roleMenu)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
