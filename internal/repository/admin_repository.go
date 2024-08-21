package repository

import (
	"context"
	"github.com/IkezawaYuki/popple/internal/domain"
	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (a *AdminRepository) FindAll(ctx context.Context) ([]domain.AdminDto, error) {
	var admins []domain.AdminDto
	err := a.db.WithContext(ctx).Find(&admins).Error
	return admins, err
}

func (a *AdminRepository) FindById(ctx context.Context, id uint64) (*domain.AdminDto, error) {
	var admin domain.AdminDto
	err := a.db.WithContext(ctx).First(&admin, id).Error
	return &admin, err
}

func (a *AdminRepository) FindByEmail(ctx context.Context, email string) (*domain.AdminDto, error) {
	var admin domain.AdminDto
	err := a.db.WithContext(ctx).First(&admin, "email = ?", email).Error
	return &admin, err
}
