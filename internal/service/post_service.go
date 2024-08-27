package service

import (
	"context"
	"errors"
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"github.com/IkezawaYuki/popple/internal/repository"
	"gorm.io/gorm"
)

type PostService struct {
	postRepo *repository.PostRepository
}

func NewPostService(postRepo *repository.PostRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
	}
}

func (s *PostService) IsLinked(ctx context.Context, instagramMediaID string) (bool, error) {
	_, err := s.postRepo.FindByInstagramMediaID(ctx, instagramMediaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *PostService) Create(ctx context.Context, post *entity.Post) error {
	m := model.Post{}
	m.CustomerID = post.CustomerID
	m.InstagramMediaID = post.InstagramMediaID
	return s.postRepo.Save(ctx, &m)
}

func (s *PostService) SaveInstagramPost(ctx context.Context, customerID int, post *entity.InstagramPost) (*entity.Post, error) {
	m := model.Post{}
	m.CustomerID = customerID
	m.InstagramMediaID = post.ID
	m.InstagramLink = post.MediaURL
	err := s.postRepo.Save(ctx, &m)
	if err != nil {
		return nil, err
	}
	return &entity.Post{
		ID:               int(m.ID),
		CustomerID:       m.CustomerID,
		InstagramMediaID: m.InstagramMediaID,
		CreatedAt:        m.CreatedAt,
	}, nil
}

func (s *PostService) SaveWordpressPost(ctx context.Context, post *entity.Post) error {
	m, err := s.postRepo.FindByInstagramMediaID(ctx, post.InstagramMediaID)
	if err != nil {
		return err
	}
	m.WordpressLink.String = *post.WordpressLink
	return s.postRepo.Save(ctx, m)
}
