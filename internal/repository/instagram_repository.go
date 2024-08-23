package repository

import (
	"context"
	"errors"
	"github.com/IkezawaYuki/popple/internal/domain"
	"gorm.io/gorm"
)

type InstagramRepository struct {
	db *gorm.DB
}

func NewInstagramRepository(db *gorm.DB) *InstagramRepository {
	return &InstagramRepository{db: db}
}

func (i *InstagramRepository) FindByCustomerUUID(ctx context.Context, customerID string) ([]domain.InstagramPostDto, error) {
	var dto []domain.InstagramPostDto
	err := i.db.WithContext(ctx).
		Preload("InstagramPostMediaDto").
		Find(&dto, "customer_id = ?", customerID).Error
	return dto, err
}

func (i *InstagramRepository) FindNotYetByCustomerUUID(ctx context.Context, customerID string) ([]domain.InstagramPostDto, error) {
	var dto []domain.InstagramPostDto
	err := i.db.WithContext(ctx).
		Preload("InstagramPostMediaDto").
		Find(&dto, "customer_id = ? and post_status = 0", customerID).Error
	return dto, err
}

func (i *InstagramRepository) Save(ctx context.Context, dto domain.InstagramPostDto) error {
	err := i.db.WithContext(ctx).Save(dto).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicateKey
		}
		return err
	}
	return nil
}
