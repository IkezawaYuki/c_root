package usecase

import (
	"context"
	"fmt"
	"github.com/IkezawaYuki/c_root/internal/domain"
	"github.com/IkezawaYuki/c_root/internal/service"
)

type CustomerUsecase struct {
	customerService service.CustomerService
	authService     service.AuthService
}

func NewCustomerUsecase(customerService service.CustomerService, authService service.AuthService) CustomerUsecase {
	return CustomerUsecase{
		customerService: customerService,
		authService:     authService,
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
	adminService    service.AdminService
	authService     service.AuthService
	customerService service.CustomerService
}

func NewAdminUsecase(
	adminSrv service.AdminService,
	authSrv service.AuthService,
	customerSrv service.CustomerService,
) AdminUsecase {
	return AdminUsecase{
		adminService:    adminSrv,
		authService:     authSrv,
		customerService: customerSrv,
	}
}

func (a *AdminUsecase) RegisterCustomer(ctx context.Context, customer *domain.Customer) error {
	return a.customerService.CreateCustomer(ctx, customer)
}
