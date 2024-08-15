package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/IkezawaYuki/c_root/internal/domain"
	"github.com/IkezawaYuki/c_root/internal/repository"
	"github.com/IkezawaYuki/c_root/internal/service"
)

type CustomerUsecase struct {
	baseRepository   repository.BaseRepository
	customerService  service.CustomerService
	authService      service.AuthService
	wordpressService service.WordpressService
	instagramService service.InstagramService
	metaService      service.MetaService
}

func NewCustomerUsecase(baseRepo repository.BaseRepository,
	customerSrv service.CustomerService,
	authSrv service.AuthService,
	wordpressSrv service.WordpressService,
	instagramSrv service.InstagramService,
	metaSrv service.MetaService,
) CustomerUsecase {
	return CustomerUsecase{
		baseRepository:   baseRepo,
		customerService:  customerSrv,
		authService:      authSrv,
		wordpressService: wordpressSrv,
		instagramService: instagramSrv,
		metaService:      metaSrv,
	}
}

func (c *CustomerUsecase) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	return c.customerService.GetCustomer(ctx, id)
}

func (c *CustomerUsecase) Login(ctx context.Context, user *domain.User) (string, error) {
	customer, err := c.customerService.GetCustomerByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}
	if err := c.authService.CheckPassword(user, customer); err != nil {
		return "", fmt.Errorf("invalid password: %w", err)
	}
	return c.authService.GenerateJWTCustomer(customer)
}

func (c *CustomerUsecase) GetInstagramMedia(ctx context.Context, id string) ([]*domain.InstagramMediaDetail, error) {
	customer, err := c.customerService.GetCustomer(ctx, id)
	if err != nil {
		return nil, err
	}
	if customer.FacebookToken == nil {
		return nil, errors.New("invalid operation")
	}
	mediaList, err := c.instagramService.GetMediaList(ctx, *customer.FacebookToken)
	if err != nil {
		return nil, err
	}
	var media []*domain.InstagramMediaDetail
	for _, mediaID := range mediaList {
		detail, err := c.instagramService.GetMediaDetail(ctx, *customer.FacebookToken, mediaID)
		if err != nil {
			return nil, err
		}
		media = append(media, detail)
	}
	return media, nil
}

type AdminUsecase struct {
	baseRepository  repository.BaseRepository
	adminService    service.AdminService
	authService     service.AuthService
	customerService service.CustomerService
}

func NewAdminUsecase(
	baseRepo repository.BaseRepository,
	adminSrv service.AdminService,
	authSrv service.AuthService,
	customerSrv service.CustomerService,
) AdminUsecase {
	return AdminUsecase{
		baseRepository:  baseRepo,
		adminService:    adminSrv,
		authService:     authSrv,
		customerService: customerSrv,
	}
}

func (a *AdminUsecase) RegisterCustomer(ctx context.Context, customer *domain.Customer) error {
	return a.customerService.CreateCustomer(ctx, customer)
}
