package usecase

import (
	"context"
	"github.com/IkezawaYuki/popple/internal/domain"
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/domain/objects"
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

func (a *AdminUsecase) GetLinkHistories(ctx context.Context) ([]domain.LinkHistory, error) {
	return a.linkHistoryService.GetLinkHistories(ctx)
}

func (a *AdminUsecase) GetLinkHistoriesByCustomer(ctx context.Context, customerUUID string) ([]domain.LinkHistory, error) {
	return a.linkHistoryService.GetLinkHistoriesByCustomer(ctx, customerUUID)
}

func (a *AdminUsecase) GetLinkHistoryByUUID(ctx context.Context, uuid string) (*domain.LinkHistory, error) {
	return a.linkHistoryService.GetLinkHistoryByUUID(ctx, uuid)
}
