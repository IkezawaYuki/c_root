package repository

import (
	"context"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewInstagramWordpressRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (p *PostRepository) Save(ctx context.Context, post *model.Post) error {
	return p.db.WithContext(ctx).Save(post).Error
}

func (p *PostRepository) FindAll(ctx context.Context) ([]model.Post, error) {
	var posts []model.Post
	err := p.db.WithContext(ctx).Find(&posts).Order("customer_id").Error
	return posts, err
}

func (p *PostRepository) FindByCustomerID(ctx context.Context, customerUUID string) ([]model.Post, error) {
	var posts []model.Post
	err := p.db.WithContext(ctx).Find(&posts, "customer_id = ?", customerUUID).Error
	return posts, err
}

func (p *PostRepository) FindByUUID(ctx context.Context, uuid string) (*model.Post, error) {
	var post model.Post
	err := p.db.WithContext(ctx).Where("uuid = ?", uuid).First(&post).Error
	return &post, err
}
