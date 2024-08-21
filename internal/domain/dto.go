package domain

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type CustomerDto struct {
	UUID           string
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

func (c *CustomerDto) TableName() string {
	return "customers"
}

func (c *CustomerDto) ConvertToCustomer() *Customer {
	var customer Customer
	customer.ID = c.ID
	customer.Name = c.Name
	customer.Email = c.Email
	customer.Password = c.Password
	customer.WordpressURL = c.WordpressURL
	customer.DeleteHashFlag = c.DeleteHashFlag
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

type AdminDto struct {
	UUID     string
	Name     string
	Email    string
	Password string
	gorm.Model
}

func (AdminDto) TableName() string { return "admins" }

type InstagramWordpressDto struct {
	UUID             string
	CustomerUUID     string
	InstagramMediaID int
	InstagramLink    string
	WordpressMediaID int
	WordpressLink    string
	gorm.Model
}

func (InstagramWordpressDto) TableName() string { return "instagram_wordpress" }

type InstagramDto struct {
	UUID       string
	MediaID    int
	Caption    string
	MediaType  string
	MediaURL   string
	Permalink  string
	PostStatus int
	Timestamp  time.Time
	gorm.Model
}

func (InstagramDto) TableName() string { return "instagram_posts" }
