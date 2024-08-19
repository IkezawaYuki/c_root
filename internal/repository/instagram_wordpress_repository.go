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
	return &InstagramWordpressRepository{}
}

func (i *InstagramWordpressRepository) Save(ctx context.Context, dto domain.InstagramWordpressDto) error {
	return i.db.WithContext(ctx).Save(dto).Error
}
