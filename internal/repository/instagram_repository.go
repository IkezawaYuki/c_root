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

func (i *InstagramRepository) FindByCustomerID(ctx context.Context, customerID string) ([]domain.InstagramDto, error) {
	var dto []domain.InstagramDto
	err := i.db.WithContext(ctx).Find(&dto, "customer_id = ?", customerID).Error
	return dto, err
}

func (i *InstagramRepository) FindNotYetByCustomerID(ctx context.Context, customerID string) ([]domain.InstagramDto, error) {
	var dto []domain.InstagramDto
	err := i.db.WithContext(ctx).Find(&dto, "customer_id = ? and post_status = 0", customerID).Error
	return dto, err
}

func (i *InstagramRepository) Save(ctx context.Context, dto domain.InstagramDto) error {
	err := i.db.WithContext(ctx).Save(dto).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicateKey
		}
		return err
	}
	return nil
}
