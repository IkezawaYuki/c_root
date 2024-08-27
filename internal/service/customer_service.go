package service

import (
	"context"
	"errors"
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"github.com/IkezawaYuki/popple/internal/domain/objects"
	"github.com/IkezawaYuki/popple/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CustomerService struct {
	customerRepository *repository.CustomerRepository
	postRepository     *repository.PostRepository
}

func NewCustomerService(
	customerRepo *repository.CustomerRepository,
	postRepository *repository.PostRepository,
) *CustomerService {
	return &CustomerService{
		customerRepository: customerRepo,
		postRepository:     postRepository,
	}
}

func (s *CustomerService) FindAll(ctx context.Context) ([]entity.Customer, error) {
	postModelList, err := s.customerRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	customers := make([]entity.Customer, len(postModelList))
	for i, m := range postModelList {
		customer := convertCustomer(&m)
		customers[i] = *customer
	}
	return customers, nil
}

func (s *CustomerService) FindByID(ctx context.Context, id int) (*entity.Customer, error) {
	m, err := s.customerRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return convertCustomer(m), nil
}

func (s *CustomerService) GetCustomerByEmail(ctx context.Context, email string) (*entity.Customer, error) {
	m, err := s.customerRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return convertCustomer(m), nil
}

func convertCustomer(m *model.Customer) *entity.Customer {
	var customer entity.Customer

	customer.ID = int(m.ID)
	customer.Name = m.Name
	customer.Password = m.Password
	customer.Email = m.Email
	customer.WordpressURL = m.WordpressURL

	if m.FacebookToken.Valid {
		customer.FacebookToken = &m.FacebookToken.String
	}
	if m.InstagramID.Valid {
		customer.InstagramID = &m.InstagramID.String
	}
	if m.InstagramName.Valid {
		customer.InstagramName = &m.InstagramName.String
	}
	customer.DeleteHashFlag = m.DeleteHashFlag
	return &customer
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customer *entity.Customer) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := s.customerRepository.Save(ctx, &model.Customer{
		Name:           customer.Name,
		Email:          customer.Email,
		Password:       string(passwordHash),
		DeleteHashFlag: 0,
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return objects.ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, id string) error {
	panic("implement me")
}

//func (s *CustomerService) GetInstagramPostNotYet(ctx context.Context, customerID string) ([]domain.InstagramPost, error) {
//	records, err := s.instagramRepository.FindNotYetByCustomerUUID(ctx, customerID)
//	if err != nil {
//		return nil, err
//	}
//	posts := make([]domain.InstagramPost, len(records))
//	for i, record := range records {
//		medias := make([]domain.Media, len(record.InstagramPostMediaDto))
//		for j, child := range record.InstagramPostMediaDto {
//			medias[j] = domain.Media{
//				ID:   child.MediaID,
//				Url:  child.MediaURL,
//				Type: child.MediaType,
//			}
//		}
//		posts[i] = domain.InstagramPost{
//			ID:         record.ID,
//			Caption:    record.Caption,
//			MediaType:  record.Caption,
//			MediaURL:   record.MediaURL,
//			PostStatus: domain.PostStatus(record.PostStatus),
//			Timestamp:  record.Timestamp,
//			Children:   medias,
//		}
//	}
//	return posts, nil
//}

//func (s *CustomerService) GetInstagramPost(ctx context.Context, customerID string) ([]domain.InstagramPost, error) {
//	records, err := s.instagramRepository.FindByCustomerUUID(ctx, customerID)
//	if err != nil {
//		return nil, err
//	}
//	posts := make([]domain.InstagramPost, len(records))
//	for i, record := range records {
//		medias := make([]domain.Media, len(record.InstagramPostMediaDto))
//		for j, child := range record.InstagramPostMediaDto {
//			medias[j] = domain.Media{
//				ID:   child.MediaID,
//				Url:  child.MediaURL,
//				Type: child.MediaType,
//			}
//		}
//		posts[i] = domain.InstagramPost{
//			ID:         record.ID,
//			Caption:    record.Caption,
//			MediaType:  record.Caption,
//			MediaURL:   record.MediaURL,
//			PostStatus: domain.PostStatus(record.PostStatus),
//			Timestamp:  record.Timestamp,
//			Children:   medias,
//		}
//	}
//	return posts, nil
//}

//func (s *CustomerService) SaveInstagramPost(ctx context.Context, instagramPost *entity.InstagramMediaDetail, instagramPostMedia []domain.Media, startDate *time.Time) error {
//	timestamp, err := time.Parse("2006-01-02T15:04:05-0700", instagramPost.Timestamp)
//	if err != nil {
//		return err
//	}
//	if startDate == nil {
//		return errors.New("startDate is required")
//	}
//	status := entity.NotYet
//	if startDate.Before(timestamp) {
//		status = entity.Linked
//	}
//	mediaDto := make([]model.InstagramPostMediaDto, len(instagramPostMedia))
//	for i, media := range instagramPostMedia {
//		mediaDto[i] = model.InstagramPostMediaDto{
//			UUID:      uuid.NewString(),
//			MediaID:   media.ID,
//			MediaURL:  media.Url,
//			MediaType: media.Type,
//		}
//	}
//	err = s.instagramRepository.Save(ctx, model.InstagramPostDto{
//		UUID:                  uuid.NewString(),
//		Caption:               instagramPost.Caption,
//		MediaType:             instagramPost.MediaType,
//		MediaURL:              instagramPost.MediaURL,
//		Permalink:             instagramPost.Permalink,
//		PostStatus:            int(status),
//		Timestamp:             timestamp,
//		InstagramPostMediaDto: mediaDto,
//	})
//	if err != nil {
//		if errors.Is(err, objects.ErrDuplicateKey) {
//			return nil
//		}
//		return err
//	}
//	return nil
//}
//
//func (s *CustomerService) CreateInstagramWordpress(ctx context.Context, instagramLink, wordpressLink string) error {
//	return s.instagramWordpressRepository.Save(ctx, model.InstagramWordpressDto{
//		UUID:          uuid.New().String(),
//		WordpressLink: wordpressLink,
//		InstagramLink: instagramLink,
//	})
//}
