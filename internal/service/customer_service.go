package service

import (
	"context"
	"errors"
	"github.com/IkezawaYuki/popple/internal/domain"
	"github.com/IkezawaYuki/popple/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type CustomerService struct {
	customerRepository           *repository.CustomerRepository
	instagramRepository          *repository.InstagramRepository
	instagramWordpressRepository *repository.InstagramWordpressRepository
}

func NewCustomerService(
	customerRepo *repository.CustomerRepository,
	instagramRepository *repository.InstagramRepository,
	instagramWordpressRepository *repository.InstagramWordpressRepository,
) *CustomerService {
	return &CustomerService{
		customerRepository:           customerRepo,
		instagramRepository:          instagramRepository,
		instagramWordpressRepository: instagramWordpressRepository,
	}
}

func (s *CustomerService) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	return s.customerRepository.FindByID(ctx, id)
}

func (s *CustomerService) GetCustomerByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	return s.customerRepository.FindByEmail(ctx, email)
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customer *domain.Customer) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	customer.ID = uuid.New().String()
	if err := s.customerRepository.Save(ctx, &domain.CustomerDto{
		ID:             customer.ID,
		Name:           customer.Name,
		Email:          customer.Email,
		Password:       string(passwordHash),
		DeleteHashFlag: 0,
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, id string) error {
	panic("implement me")
}

func (s *CustomerService) GetInstagramPostNotYet(ctx context.Context, customerID string) ([]domain.InstagramPost, error) {
	records, err := s.instagramRepository.FindNotYetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	posts := make([]domain.InstagramPost, len(records))
	for i, record := range records {
		posts[i] = domain.InstagramPost{
			ID:         record.ID,
			Caption:    record.Caption,
			MediaType:  record.Caption,
			MediaURL:   record.MediaURL,
			PostStatus: domain.PostStatus(record.PostStatus),
			Timestamp:  record.Timestamp,
		}
	}
	return posts, nil
}

func (s *CustomerService) GetInstagramPost(ctx context.Context, customerID string) ([]domain.InstagramPost, error) {
	records, err := s.instagramRepository.FindByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	posts := make([]domain.InstagramPost, len(records))
	for i, record := range records {
		posts[i] = domain.InstagramPost{
			ID:         record.ID,
			Caption:    record.Caption,
			MediaType:  record.Caption,
			MediaURL:   record.MediaURL,
			PostStatus: domain.PostStatus(record.PostStatus),
			Timestamp:  record.Timestamp,
		}
	}
	return posts, nil
}

func (s *CustomerService) SaveInstagramPost(ctx context.Context, instagramPost *domain.InstagramMediaDetail, startDate *time.Time) error {
	timestamp, err := time.Parse("2006-01-02T15:04:05-0700", instagramPost.Timestamp)
	if err != nil {
		return err
	}
	if startDate == nil {
		return errors.New("startDate is required")
	}
	status := domain.NotYet
	if startDate.Before(timestamp) {
		status = domain.Linked
	}
	err = s.instagramRepository.Save(ctx, domain.InstagramDto{
		ID:         instagramPost.ID,
		Caption:    instagramPost.Caption,
		MediaType:  instagramPost.MediaType,
		MediaURL:   instagramPost.MediaURL,
		Permalink:  instagramPost.Permalink,
		PostStatus: int(status),
		Timestamp:  timestamp,
	})
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateKey) {
			return nil
		}
		return err
	}
	return nil
}

func (s *CustomerService) CreateInstagramWordpress(ctx context.Context, instagramLink, wordpressLink string) error {
	return s.instagramWordpressRepository.Save(ctx, domain.InstagramWordpressDto{
		UUID:          uuid.New().String(),
		WordpressLink: wordpressLink,
		InstagramLink: instagramLink,
	})
}
