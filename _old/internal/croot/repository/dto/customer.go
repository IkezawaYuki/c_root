package dto

import (
	"database/sql"
	"github.com/IkezawaYuki/popple/internal/croot/domain"
	"gorm.io/gorm"
)

type Customer struct {
	ID             string         `gorm:"column:id"`
	Name           string         `gorm:"column:name"`
	EMail          string         `gorm:"column:email"`
	Password       string         `gorm:"column:password"`
	FacebookToken  sql.NullString `gorm:"column:facebook_token"`
	StartDate      sql.NullTime   `gorm:"column:start_date"`
	InstagramID    sql.NullString `gorm:"column:instagram_id"`
	InstagramName  sql.NullString `gorm:"column:instagram_name"`
	DeleteHashFlag int            `gorm:"column:delete_hash_flag"`
	gorm.Model
}

func (Customer) TableName() string {
	return "customers"
}

func ConvertCustomer(c Customer) *domain.Customer {
	var customer = domain.Customer{
		ID:             c.ID,
		Name:           c.Name,
		Email:          c.EMail,
		DeleteHashFlag: c.DeleteHashFlag,
	}
	if c.FacebookToken.Valid {
		customer.FacebookToken = &c.FacebookToken.String
	}
	if c.StartDate.Valid {
		customer.StartDate = &c.StartDate.Time
	}
	if c.InstagramID.Valid {
		customer.InstagramID = &c.InstagramID.String
	}
	if c.InstagramName.Valid {
		customer.InstagramName = &c.InstagramName.String
	}
	return &customer
}
