package service

import (
	"context"
	"github.com/IkezawaYuki/popple/internal/domain"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"github.com/IkezawaYuki/popple/internal/repository"
)

type LinkHistoryService struct {
	instaWordpressRepo *repository.InstagramWordpressRepository
}

func NewLinkHistoryService(instaWordpressRepo *repository.InstagramWordpressRepository) *LinkHistoryService {
	return &LinkHistoryService{
		instaWordpressRepo: instaWordpressRepo,
	}
}

func (s *LinkHistoryService) IsLinked(ctx context.Context, link *domain.LinkHistory) (bool, error) {

}

func (s *LinkHistoryService) Create(ctx context.Context, link *domain.LinkHistory) error {
	return s.instaWordpressRepo.Save(ctx, model.InstagramWordpressDto{
		CustomerID:       "",
		InstagramMediaID: "",
		InstagramLink:    "",
		WordpressMediaID: "",
		WordpressLink:    "",
	})
}

func (s *LinkHistoryService) GetLinkHistories(ctx context.Context) ([]domain.LinkHistory, error) {
	dtoList, err := s.instaWordpressRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	histories := make([]domain.LinkHistory, len(dtoList))
	for i, dto := range dtoList {
		histories[i] = domain.LinkHistory{
			ID:               int(dto.ID),
			UUID:             dto.UUID,
			InstagramMediaID: dto.InstagramMediaID,
			InstagramLink:    dto.InstagramLink,
			WordpressMediaID: dto.WordpressMediaID,
			WordpressLink:    dto.WordpressLink,
			CreateAt:         dto.CreatedAt,
		}
	}
	return histories, nil
}

func (s *LinkHistoryService) GetLinkHistoriesByCustomer(ctx context.Context, customerUUID string) ([]domain.LinkHistory, error) {
	dtoList, err := s.instaWordpressRepo.FindByCustomerID(ctx, customerUUID)
	if err != nil {
		return nil, err
	}
	histories := make([]domain.LinkHistory, len(dtoList))
	for i, dto := range dtoList {
		histories[i] = domain.LinkHistory{
			ID:               int(dto.ID),
			UUID:             dto.UUID,
			InstagramMediaID: dto.InstagramMediaID,
			InstagramLink:    dto.InstagramLink,
			WordpressMediaID: dto.WordpressMediaID,
			WordpressLink:    dto.WordpressLink,
			CreateAt:         dto.CreatedAt,
		}
	}
	return histories, nil
}

func (s *LinkHistoryService) GetLinkHistoryByUUID(ctx context.Context, uuid string) (*domain.LinkHistory, error) {
	dto, err := s.instaWordpressRepo.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return &domain.LinkHistory{
		ID:               int(dto.ID),
		UUID:             dto.UUID,
		InstagramMediaID: dto.InstagramMediaID,
		InstagramLink:    dto.InstagramLink,
		WordpressMediaID: dto.WordpressMediaID,
		WordpressLink:    dto.WordpressLink,
		CreateAt:         dto.CreatedAt,
	}, nil
}
