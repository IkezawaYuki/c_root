package usecase

import (
	"context"
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/domain/objects"
	"github.com/IkezawaYuki/popple/internal/repository"
	"github.com/IkezawaYuki/popple/internal/service"
)

type AdminUsecase struct {
	baseRepository  *repository.BaseRepository
	adminService    *service.AdminService
	authService     *service.AuthService
	customerService *service.CustomerService
	postService     *service.PostService
	customerUsecase *CustomerUsecase
}

func NewAdminUsecase(
	baseRepo *repository.BaseRepository,
	adminSrv *service.AdminService,
	authSrv *service.AuthService,
	customerSrv *service.CustomerService,
	customerUsecase *CustomerUsecase,
) *AdminUsecase {
	return &AdminUsecase{
		baseRepository:  baseRepo,
		adminService:    adminSrv,
		authService:     authSrv,
		customerService: customerSrv,
		customerUsecase: customerUsecase,
	}
}

func (a *AdminUsecase) RegisterCustomer(ctx context.Context, customer *entity.Customer) error {
	return a.customerService.CreateCustomer(ctx, customer)
}

func (a *AdminUsecase) Login(ctx context.Context, user *entity.User) (string, error) {
	customer, err := a.adminService.FindByEmail(ctx, user.Email)
	if err != nil {
		return "", objects.ErrNotFound
	}
	if err := a.authService.CheckPassword(user, customer.Password); err != nil {
		return "", err
	}
	return a.authService.GenerateJWTAdmin(customer)
}

func (a *AdminUsecase) GetCustomers(ctx context.Context) ([]entity.Customer, error) {
	return a.customerService.FindAll(ctx)
}

func (a *AdminUsecase) GetAdmins(ctx context.Context) ([]entity.Admin, error) {
	return a.adminService.FindAll(ctx)
}
