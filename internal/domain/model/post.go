package model

import (
	"database/sql"
	"gorm.io/gorm"
)

type Post struct {
	CustomerID       string `gorm:"not null"`
	InstagramMediaID string `gorm:"not null"`
	InstagramLink    string `gorm:"not null"`
	WordpressMediaID sql.NullString
	WordpressLink    sql.NullString
	gorm.Model
}

func (Post) TableName() string { return "posts" }
