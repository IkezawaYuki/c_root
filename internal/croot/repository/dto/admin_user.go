package dto

import (
	"github.com/IkezawaYuki/c_root/internal/croot/domain/entity"
	"gorm.io/gorm"
)

type AdminUser struct {
	ID       string `gorm:"column:id"`
	Name     string `gorm:"column:name"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
	gorm.Model
}

func (AdminUser) TableName() string { return "admin_users" }

func ConvertAdminUser(a AdminUser) *entity.AdminUser {
	return &entity.AdminUser{
		ID:       a.ID,
		Name:     a.Name,
		Email:    a.Email,
		Password: a.Password,
	}
}
