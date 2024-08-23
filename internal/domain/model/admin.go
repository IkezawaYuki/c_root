package model

import "gorm.io/gorm"

type Admin struct {
	Name     string `gorm:"not null"`
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	gorm.Model
}

func (Admin) TableName() string { return "admins" }
