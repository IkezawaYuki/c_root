package model

import (
	"database/sql"
	"gorm.io/gorm"
)

type Customer struct {
	Name           string `gorm:"not null"`
	Email          string `gorm:"unique;not null"`
	Password       string `gorm:"not null"`
	WordpressURL   string `gorm:"unique;not null"`
	FacebookToken  sql.NullString
	StartDate      sql.NullTime
	InstagramID    sql.NullString
	InstagramName  sql.NullString
	DeleteHashFlag int
	gorm.Model
}

func (c *Customer) TableName() string {
	return "customers"
}
