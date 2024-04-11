package repository

import (
	"go-authorization/errors"
	"go-authorization/lib"
	"go-authorization/models"
	"gorm.io/gorm"
)

type IRoleRepository interface {
	WithTrx(trxHandle *gorm.DB) IRoleRepository
	Query(param *models.RoleQueryParam) (*models.RoleQueryResult, error)
	Get(id string) (*models.Role, error)
	Create(role *models.Role) error
	Update(id string, role *models.Role) error
	Delete(id string) error
	UpdateStatus(id string, status int) error
}

// roleRepository database structure
type roleRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db lib.Database, logger lib.Logger) IRoleRepository {
	return roleRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a roleRepository) WithTrx(trxHandle *gorm.DB) IRoleRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a roleRepository) Query(param *models.RoleQueryParam) (*models.RoleQueryResult, error) {
	db := a.db.ORM.Model(&models.Role{})

	if v := param.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	if v := param.Name; v != "" {
		db = db.Where("name=?", v)
	}

	if v := param.UserID; v != "" {
		subQuery := a.db.ORM.Model(&models.UserRole{}).
			Where("user_id=?", v).
			Select("role_id")

		db = db.Where("id IN (?)", subQuery)
	}

	if v := param.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("name LIKE ? OR remark LIKE ?", v, v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.Roles, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.RoleQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a roleRepository) Get(id string) (*models.Role, error) {
	role := new(models.Role)

	if ok, err := QueryOne(a.db.ORM.Model(role).Where("id=?", id), role); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return role, nil
}

func (a roleRepository) Create(role *models.Role) error {
	result := a.db.ORM.Model(role).Create(role)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a roleRepository) Update(id string, role *models.Role) error {
	result := a.db.ORM.Model(role).Where("id=?", id).Updates(role)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a roleRepository) Delete(id string) error {
	role := new(models.Role)

	result := a.db.ORM.Model(role).Where("id=?", id).Delete(role)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a roleRepository) UpdateStatus(id string, status int) error {
	role := new(models.Role)

	result := a.db.ORM.Model(role).Where("id=?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
