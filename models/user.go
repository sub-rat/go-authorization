package models

import (
	"go-authorization/models/database"
	"go-authorization/models/dto"
)

// User
// Status - 1: Enable 0: Disable
type User struct {
	database.Model
	ID        string    `gorm:"column:id;size:36;index;not null;" json:"id"`
	Username  string    `gorm:"column:username;size:64;not null;index;" json:"username" validate:"required"`
	FullName  string    `gorm:"column:full_name;size:64;not null;" json:"full_name" validate:"required"`
	Password  string    `gorm:"column:password;not null;" json:"password" json:"phone"`
	Email     string    `gorm:"column:email;default:'';" json:"email"`
	Phone     string    `gorm:"column:phone;default:'';" json:"phone"`
	Status    int       `gorm:"column:status;not null;default:0;" json:"status" validate:"required,max=1,min=-1"`
	CreatedBy string    `gorm:"column:created_by;not null;" json:"created_by"`
	UserRoles UserRoles `gorm:"-" json:"user_roles"`
}

type Users []*User

type UserInfo struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Roles    Roles  `json:"roles"`
}

type UserQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	QueryPassword bool
	Username      string   `query:"username"`
	FullName      string   `query:"full_name"`
	QueryValue    string   `query:"query_value"`
	Status        int      `query:"status" validate:"max=1,min=-1"`
	RoleIDs       []string `query:"-"`
}

type UserQueryResult struct {
	List       Users           `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a *User) CleanSecure() *User {
	a.Password = ""
	return a
}

func (a Users) ToIDs() []string {
	ids := make([]string, len(a))
	for i, item := range a {
		ids[i] = item.ID
	}
	return ids
}
