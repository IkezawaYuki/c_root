package usecase

import (
	"context"
	"github.com/IkezawaYuki/c_root/internal/domain"
	"github.com/IkezawaYuki/c_root/internal/repository"
	"github.com/IkezawaYuki/c_root/internal/service"
)

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
