package repository

import (
	"go-authorization/errors"
	"go-authorization/lib"
	"go-authorization/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	WithTrx(trxHandle *gorm.DB) IUserRepository
	Query(param *models.UserQueryParam) (*models.UserQueryResult, error)
	Get(id string) (*models.User, error)
	Create(user *models.User) error
	Update(id string, user *models.User) error
	Delete(id string) error
	UpdateStatus(id string, status int) error
	UpdatePassword(id, password string) error
}

// userRepository database structure
type userRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewUserRepository creates a new user repository
func NewUserRepository(db lib.Database, logger lib.Logger) IUserRepository {
	return userRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a userRepository) WithTrx(trxHandle *gorm.DB) IUserRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

// Query  gets all users
func (a userRepository) Query(param *models.UserQueryParam) (*models.UserQueryResult, error) {
	db := a.db.ORM.Model(&models.User{})

	if v := param.QueryPassword; !v {
		db = db.Omit("password")
	}

	if v := param.Username; v != "" {
		db = db.Where("username = (?)", v)
	}

	if v := param.FullName; v != "" {
		db = db.Where("full_name = (?)", v)
	}

	if v := param.Status; v != 0 {
		db = db.Where("status = (?)", v)
	}

	if v := param.RoleIDs; len(v) > 0 {
		subQuery := a.db.ORM.Model(&models.UserRole{}).
			Select("user_id").
			Where("role_id IN (?)", v)

		db = db.Where("id IN (?)", subQuery)
	}

	if v := param.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("username LIKE ? OR full_name LIKE ? OR phone LIKE ? OR email LIKE ?", v, v, v, v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.Users, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.UserQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a userRepository) Get(id string) (*models.User, error) {
	user := new(models.User)

	if ok, err := QueryOne(a.db.ORM.Model(user).Where("id=?", id), user); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return user, nil
}

func (a userRepository) Create(user *models.User) error {
	result := a.db.ORM.Model(user).Create(user)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a userRepository) Update(id string, user *models.User) error {
	result := a.db.ORM.Model(user).Where("id=?", id).Updates(user)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a userRepository) Delete(id string) error {
	user := new(models.User)

	result := a.db.ORM.Model(user).Where("id=?", id).Delete(user)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a userRepository) UpdateStatus(id string, status int) error {
	user := new(models.User)

	result := a.db.ORM.Model(user).Where("id=?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a userRepository) UpdatePassword(id, password string) error {
	user := new(models.User)

	result := a.db.ORM.Model(user).Where("id=?", id).Update("password", password)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
