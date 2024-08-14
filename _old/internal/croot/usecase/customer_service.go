package usecase

import (
	"context"
	"github.com/IkezawaYuki/c_root/internal/croot/domain"
	"github.com/IkezawaYuki/c_root/internal/croot/domain/crooterrors"
	"github.com/IkezawaYuki/c_root/internal/croot/interfaces"
)

type CustomerService interface {
	GetInstagram(ctx context.Context, customerId string) ([]domain.InstagramMediaDetail, error)
}

type customerService struct {
	repository    interfaces.Repository
	customerRepo  domain.CustomerRepository
	instagramRepo domain.InstagramPostRepository
	wordpressRepo domain.WordpressRepository
}

func NewCustomerService(customerRepo domain.CustomerRepository, instagramRepo domain.InstagramPostRepository) CustomerService {
	return &customerService{
		customerRepo:  customerRepo,
		instagramRepo: instagramRepo,
	}
}

func (c *customerService) GetInstagram(ctx context.Context, customerId string) ([]domain.InstagramMediaDetail, error) {
	tx, err := c.repository.Begin(ctx)
	if err != nil {
		return nil, err
	}

	customer, err := c.customerRepo.FindByIDWithTX(ctx, customerId, tx)
	if err != nil {
		return nil, err
	}

	if customer.FacebookToken == nil {
		return nil, crooterrors.UnauthorizedError
	}

	instagramList, err := c.instagramRepo.FindAll(ctx, *customer.FacebookToken)
	if err != nil {
		return nil, err
	}

	result := make([]domain.InstagramMediaDetail, len(instagramList.Media.Data))
	for i, data := range instagramList.Media.Data {
		detail, err := c.instagramRepo.FindDetail(ctx, *customer.FacebookToken, data.ID)
		if err != nil {
			return nil, err
		}
		result[i] = *detail
	}

	return result, nil
}