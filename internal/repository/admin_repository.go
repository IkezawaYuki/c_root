package repository

import (
	"context"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (a *AdminRepository) FindAll(ctx context.Context) ([]model.Admin, error) {
	var admins []model.Admin
	err := a.db.WithContext(ctx).Find(&admins).Error
	return admins, err
}

func (a *AdminRepository) FindById(ctx context.Context, id uint64) (*model.Admin, error) {
	var admin model.Admin
	err := a.db.WithContext(ctx).First(&admin, id).Error
	return &admin, err
}

func (a *AdminRepository) FindByEmail(ctx context.Context, email string) (*model.Admin, error) {
	var admin model.Admin
	err := a.db.WithContext(ctx).First(&admin, "email = ?", email).Error
	return &admin, err
}

func (a *AdminRepository) Save(ctx context.Context, admin *model.Admin) error {
	return a.db.WithContext(ctx).Save(admin).Error
}
