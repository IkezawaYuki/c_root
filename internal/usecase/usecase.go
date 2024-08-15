package usecase

import (
	"context"
	"fmt"
	"github.com/IkezawaYuki/c_root/internal/domain"
	"github.com/IkezawaYuki/c_root/internal/repository"
	"github.com/IkezawaYuki/c_root/internal/service"
)

type CustomerUsecase struct {
	baseRepository  repository.BaseRepository
	customerService service.CustomerService
	authService     service.AuthService
}

func NewCustomerUsecase(baseRepo repository.BaseRepository, customerSrv service.CustomerService, authSrv service.AuthService) CustomerUsecase {
	return CustomerUsecase{
		baseRepository:  baseRepo,
		customerService: customerSrv,
		authService:     authSrv,
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
