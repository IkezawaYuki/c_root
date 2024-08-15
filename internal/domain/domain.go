package domain

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
)

var (
	ErrAuthentication = errors.New("authentication err")
	ErrAuthorization  = errors.New("authorization err")
)

type Customer struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type CustomerDto struct {
	ID             string         `gorm:"column:id"`
	Name           string         `gorm:"column:name"`
	Email          string         `gorm:"column:email"`
	Password       string         `gorm:"column:password"`
	FacebookToken  sql.NullString `gorm:"column:facebook_token"`
	StartDate      sql.NullTime   `gorm:"column:start_date"`
	InstagramID    sql.NullString `gorm:"column:instagram_id"`
	InstagramName  sql.NullString `gorm:"column:instagram_name"`
	DeleteHashFlag int            `gorm:"column:delete_hash_flag"`
	gorm.Model
}

func (c *CustomerDto) TableName() string {
	return "customers"
}

func (c *CustomerDto) ConvertToCustomer() *Customer {
	return &Customer{
		ID:       c.ID,
		Name:     c.Name,
		Password: c.Password,
		Email:    c.Email,
	}
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
