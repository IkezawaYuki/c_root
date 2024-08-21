package usecase

import (
	"context"
	"github.com/IkezawaYuki/popple/internal/domain"
	"github.com/IkezawaYuki/popple/internal/repository"
	"github.com/IkezawaYuki/popple/internal/service"
)

type AdminUsecase struct {
	baseRepository     *repository.BaseRepository
	adminService       *service.AdminService
	authService        *service.AuthService
	customerService    *service.CustomerService
	linkHistoryService *service.LinkHistoryService
}

func NewAdminUsecase(
	baseRepo *repository.BaseRepository,
	adminSrv *service.AdminService,
	authSrv *service.AuthService,
	customerSrv *service.CustomerService,
	linkHistoryService *service.LinkHistoryService,
) *AdminUsecase {
	return &AdminUsecase{
		baseRepository:     baseRepo,
		adminService:       adminSrv,
		authService:        authSrv,
		customerService:    customerSrv,
		linkHistoryService: linkHistoryService,
	}
}

func (a *AdminUsecase) RegisterCustomer(ctx context.Context, customer *domain.Customer) error {
	return a.customerService.CreateCustomer(ctx, customer)
}

func (a *AdminUsecase) Login(ctx context.Context, user *domain.User) error {
	return nil
}

func (a *AdminUsecase) GetCustomers(ctx context.Context) ([]domain.Customer, error) {
	return a.customerService.FindAll(ctx)
}

func (a *AdminUsecase) GetAdmins(ctx context.Context) ([]domain.Admin, error) {
	return a.adminService.FindAll(ctx)
}

func (a *AdminUsecase) GetLinkHistories(ctx context.Context) ([]domain.LinkHistory, error) {
	return a.linkHistoryService.GetLinkHistories(ctx)
}

func (a *AdminUsecase) GetLinkHistoriesByCustomer(ctx context.Context, customerUUID string) ([]domain.LinkHistory, error) {
	return a.linkHistoryService.GetLinkHistoriesByCustomer(ctx, customerUUID)
}

func (a *AdminUsecase) GetLinkHistoryByUUID(ctx context.Context, uuid string) (*domain.LinkHistory, error) {
	return a.linkHistoryService.GetLinkHistoryByUUID(ctx, uuid)
}
