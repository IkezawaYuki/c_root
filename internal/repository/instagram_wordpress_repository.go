package repository

import (
	"context"
	"github.com/IkezawaYuki/popple/internal/domain"
	"gorm.io/gorm"
)

type InstagramWordpressRepository struct {
	db *gorm.DB
}

func NewInstagramWordpressRepository(db *gorm.DB) *InstagramWordpressRepository {
	return &InstagramWordpressRepository{
		db: db,
	}
}

func (i *InstagramWordpressRepository) Save(ctx context.Context, dto domain.InstagramWordpressDto) error {
	return i.db.WithContext(ctx).Save(dto).Error
}

func (i *InstagramWordpressRepository) FindAll(ctx context.Context) ([]domain.InstagramWordpressDto, error) {
	var dtoList []domain.InstagramWordpressDto
	err := i.db.WithContext(ctx).Find(&dtoList).Order("customer_id").Error
	return dtoList, err
}

func (i *InstagramWordpressRepository) FindByCustomerID(ctx context.Context, customerUUID string) ([]domain.InstagramWordpressDto, error) {
	var dtoList []domain.InstagramWordpressDto
	err := i.db.WithContext(ctx).Find(&dtoList, "customer_id = ?", customerUUID).Error
	return dtoList, err
}

func (i *InstagramWordpressRepository) FindByUUID(ctx context.Context, uuid string) (domain.InstagramWordpressDto, error) {
	var dto domain.InstagramWordpressDto
	err := i.db.WithContext(ctx).Where("uuid = ?", uuid).First(&dto).Error
	return dto, err
}
