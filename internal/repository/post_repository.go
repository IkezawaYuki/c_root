package repository

import (
	"context"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
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

func (p *PostRepository) FindByID(ctx context.Context, id string) (*model.Post, error) {
	var post model.Post
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&post).Error
	return &post, err
}

func (p *PostRepository) FindByInstagramMediaID(ctx context.Context, InstagramMediaID string) (*model.Post, error) {
	var post model.Post
	err := p.db.WithContext(ctx).First(&post, "instagram_id = ?", InstagramMediaID).Error
	return &post, err
}
