package repository

import (
	"context"
	"errors"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"github.com/IkezawaYuki/popple/internal/domain/objects"
	"gorm.io/gorm"
)

type InstagramRepository struct {
	db *gorm.DB
}

func NewInstagramRepository(db *gorm.DB) *InstagramRepository {
	return &InstagramRepository{db: db}
}

func (i *InstagramRepository) FindByCustomerUUID(ctx context.Context, customerID string) ([]model.InstagramPostDto, error) {
	var dto []model.InstagramPostDto
	err := i.db.WithContext(ctx).
		Preload("InstagramPostMediaDto").
		Find(&dto, "customer_id = ?", customerID).Error
	return dto, err
}

func (i *InstagramRepository) Save(ctx context.Context, dto model.InstagramPostDto) error {
	err := i.db.WithContext(ctx).Save(dto).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return objects.ErrDuplicateKey
		}
		return err
	}
	return nil
}
